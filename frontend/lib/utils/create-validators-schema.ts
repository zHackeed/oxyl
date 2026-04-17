import * as Yup from 'yup';

export const createValidationScheme = <T extends object>(schema: Yup.ObjectSchema<T>) => schema;
