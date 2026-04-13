import * as Yup from 'yup';

export const createValidationScheme = <T extends Object>(schema: Yup.ObjectSchema<T>) => schema;