import { fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import { RootState } from '../store';

export const BASE_URL = process.env.API_BASE_URL;

const baseQuery = fetchBaseQuery({
  baseUrl: BASE_URL,
  prepareHeaders: (headers, { getState }) => {
    const token = (getState() as RootState)?.auth?.token ?? '';

    if (token) {
      headers.set('authorization', `Bearer ${token}`);
    }

    return headers;
  },
});

export default baseQuery;
