import { APIGatewayProxyResult } from "aws-lambda";
import express, { Request, Response } from "express";
import * as bodyParser from "body-parser";
import morgan from "morgan";

import {
    RawRouteHandler, responseHandler, parseFunctionArguments, wrapRouteHandlerArgs,
} from "./controller_builder";

function createRouteHandler(handler: RawRouteHandler) {
    return async (req: Request, res: Response) => {
        // Grab the request body.
        const body = req.body;

        // Parse the args of the function. The first arg in each handler
        // should be the for the `response` object for handling the response.
        const handlerArgs = parseFunctionArguments(handler);
        const wrappableHandlerArgs = handlerArgs.slice(1, handlerArgs.length);

        // Wrap the args for validation within the handler.
        const params = wrapRouteHandlerArgs(wrappableHandlerArgs, body);

        // Run the handler to get the response data.
        let responseData: APIGatewayProxyResult
        try {
            responseData = await handler(responseHandler, ...params);
        } catch (e) {
            console.log(e);
            responseData = responseHandler.sendWithStatusCode(500, e);
        }

        // Return the response data.
        res.status(responseData.statusCode);
        res.send(JSON.parse(responseData.body));
    };
}

export interface LocalServerRoute {
    path: string;
    method: "get" | "put" | "post" | "delete";
    handler: RawRouteHandler;
}

export function createLocalServer(routes: LocalServerRoute[]): express.Application {
    // Create the application.
    const app = express();

    // Set the server helpers.
    app.use(bodyParser.json());
    app.use(bodyParser.urlencoded({ extended: true }));
    app.use(morgan("dev"));

    for (let i = 0; i < routes.length; i++) {
        const route = routes[i];
        app[route.method](route.path, createRouteHandler(route.handler));
    }

    return app;
}
