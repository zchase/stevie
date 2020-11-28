// GenericConstructor is used for type checking classes.
export type GenericConstructor<T> = Function & { prototype: T };

/**
 * castType checks that a given value is a of a certain class.
 *
 * @param type The type you want to check against.
 * @param name The name of the class you are validating for. This is used for logging.
 * @param value The vaue you want to validate.
 */
export function castType<T>(name: string, value: any, type: GenericConstructor<T>): T {
  // Check that the value exists.
  if ((typeof value === "undefined") || (value === null)) {
    throw new Error(`Missing parameter [${name}].`)
  }

  // Check that the value is an object.
  const objCheck = typeof value === "object";
  if (!objCheck) {
      throw new TypeError(`Expected value to be an object [${name}].`);
  }

  // Check that the object is instance of the extpected type.
  if (!((value as Object) instanceof type)) {
      throw new TypeError(`Expect value [${name}] to be type [${type}].`);
  }

  // Return the validated value.
  return <T>value;
}

// TypeAssertionString are the currently support types of the assertType function.
type TypeAssertionString = "string" | "number";

/**
 * assertType checks that a value is of a given type.
 *
 * @param typeValueString The string value of the type to validate against.
 * @param value The value to validate.
 */
export function assertType<T>(typeValueString: TypeAssertionString, name: string, value: any): T {
  // Check that the value exists.
  if ((typeof value === "undefined") || (value === null)) {
      throw new Error(`Missing parameter [${name}].`)
  }

  // Validate that the value is of the expected type.
  const valueType = typeof value;
  if (valueType != typeValueString) {
      throw new TypeError(`Expect value [${name}] to be type [${typeValueString}] but got [${valueType}].`);
  }

  return <T>value;
}

// ServerArgOptionalMethods contains the types for the
// available optional methods of the ServerArgs
interface ServerArgOptionalMethods {
    cast<T>(type: GenericConstructor<T>): T | undefined;
    string(): string | undefined;
    int(): number | undefined;
}

// ServerArgObject contain the values of a wrapped server request
// arg.
export interface ServerArgObject {
    name: string;
    value: any;
    cast<T>(type: GenericConstructor<T>): T;
    string(): string;
    int(): number;
    optional: ServerArgOptionalMethods;
}

/**
 * ServerArg is a factory for wrapping incoming server requests in helper
 * functions for validating request bodies.
 *
 * @param name The name of the arg
 * @param value The value of the arg
 */
export function ServerArg(name: string, value: any): ServerArgObject {
    // Does the value exist?
    const valueExists = (typeof value !== "undefined") && (value !== null);

    // Methods for checking the type of a server arg.
    const cast = <T>(expectedType: GenericConstructor<T>): T => castType<T>(name, value, expectedType);
    const string = (): string => assertType<string>("string", name, value);
    const int = (): number => assertType<number>("number", name, value);

    // ServerArgs can be made optional by using the optional methods available.
    // Optional args can return undefined if they do not exists so use at your own peril.
    const optional: ServerArgOptionalMethods = {
        cast: <T>(type: GenericConstructor<T>) => valueExists ? cast<T>(type) : undefined,
        string: () => valueExists ? string() : undefined,
        int: () => valueExists ? int() : undefined,
    };

    return { name, value, cast, string, int, optional };
}
