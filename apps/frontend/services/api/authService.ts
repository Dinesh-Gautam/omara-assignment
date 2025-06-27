import axios from "../../lib/axios";

/**
 * Signs up a user.
 * @param token - The Firebase authentication token.
 * @returns The response data from the server.
 */
export const signUp = async (token: string) => {
  const response = await axios.post("/api/auth/signup", { token });
  return response.data;
};
