import { APIGatewayProxyEvent, APIGatewayProxyResult } from "aws-lambda";
import { ServerArg, ServerArgObject } from "./server_args";

/**
 * parseFunctionArguments parses the arguments of a function and
 * returns a string array of the arg names. The current implementation
 * only supports parsing functions written with '=>' syntax.
 *
 * Valid:
 *   const function = (arg: string) => {};
 *
 *   const obj = {
 *     method: (argName: string) => {};
 *   };
 *
 * Invalid:
 *   function namedFunction(arg: string) {};
 *
 * @param func The function you want to parse the arguments from.
 */
export function parseFunctionArguments(func: Function): string[] {
  // REGEX for identifying the args of a function in a string.
  const ARGS = /^\s*[^\(]*\(\s*([^\)]*)\)\s?=>/m;

  // REGEX for splitting function args.
  const ARG_SPLIT = /,/;

  // REGEX for valid function arg names.
  const ARG = /^\s*(_?)(\S+?)\1\s*$/;

  // REGEX for stripping comments.
  const STRIP_COMMENTS = /((\/\/.*$)|(\/\*[\s\S]*?\*\/))/mg;

  // If func is a not a function the return an empty array.
  if (typeof func !== "function") {
    return [];
  }

  // Strip the comments and identify the args of a function.
  const argNames: string[] = [];
  const funcText = func.toString().replace(STRIP_COMMENTS, '');
  const argMatches = funcText.match(ARGS) || [];

  // If there are no args in a fucntion then return an empty array.
  if (argMatches.length < 2) {
      return [];
  }

  // Loop through the matched args and push the arg name into
  // the results araary.
  const argMatch = argMatches[1].split(ARG_SPLIT);
  for (let i = 0; i < argMatch.length; i++) {
    const arg = argMatch[i];
    arg.replace(ARG, (all: any, underscore: any, name: string) => {
      argNames.push(name);
      return arg;
    });
  }

  return argNames;
};

// RawRouteHandler is the handler written by the person creating the endpoint.
export interface RawRouteHandler {
    (res: ResponseHandler, ...args: ServerArgObject[]): Promise<APIGatewayProxyResult>;
}

// RouteHandler is the handler for an API Gateway Event.
interface RouteHandler {
    (event: APIGatewayProxyEvent): Promise<APIGatewayProxyResult>;
}

// ResponseHandler contains the methods for handling an API Gateway response.
export interface ResponseHandler {
    send(data: any): APIGatewayProxyResult;
    sendWithStatusCode(code: number, data: any): APIGatewayProxyResult;
}

// createRequestResponse is a helper function for creating a valid API Gateway result.
function createRequestResponse(statusCode: number, body: string): APIGatewayProxyResult {
    return { statusCode, body, headers: {
        "Content-Type": "text/plain",
        "Access-Control-Allow-Origin": "*",
    }};
}

// createResponseObject is a helper function for creating a valid API Gateway result from a value
// that needs to be stringified.
export function createResponseObject(code: number, data: any): APIGatewayProxyResult {
    return createRequestResponse(code, JSON.stringify(data));
}

// wrapRouteHandlerArgs is a function for wrapping the args of RouteHandler as a ServerArg.
export function wrapRouteHandlerArgs(args: string[], body: any): ServerArgObject[] {
    const result: ServerArgObject[] = [];

    for (let i = 0; i < args.length; i++) {
        const argKey = args[i];
        const serverArg = ServerArg(argKey, body[argKey]);
        result.push(serverArg);
    }

    return result;
}

export const responseHandler: ResponseHandler = {
    send: (data: any) => createResponseObject(200, data),
    sendWithStatusCode: (code: number, data: any) => createResponseObject(code, data),
}

// wrapRouteHandler is a function for wrapping a route handler with common checks
// and body handling. This helps abstract away a bunch of boilerplate stuff.
export function createRouteHandler(handler: RawRouteHandler): RouteHandler {
    return async (event: APIGatewayProxyEvent): Promise<APIGatewayProxyResult> => {
        // Handle the event body.
        let body = {};
        if (event.isBase64Encoded) {
            body = event?.body ? JSON.parse(Buffer.from(event.body, "base64").toString()) : {};
        } else if (typeof event?.body === "string") {
            body = JSON.parse(event.body);
        }

        // Add the query params to the body.
        const bodyWithQuery = Object.assign({}, body, event.queryStringParameters);

        // Parse the args of the function. The first arg in each handler
        // should be the for the `response` object for handling the response.
        const handlerArgs = parseFunctionArguments(handler);
        const wrappableHandlerArgs = handlerArgs.slice(1, handlerArgs.length);

        // Wrap the args for validation within the handler.
        const params = wrapRouteHandlerArgs(wrappableHandlerArgs, bodyWithQuery);

        // Try the handler and return a 500 error for any errors.
        try {
            return await handler(responseHandler, ...params);
        } catch (e) {
            console.log(e);
            return createResponseObject(500, { ok: false, message: e.message });
        }
    };
}
