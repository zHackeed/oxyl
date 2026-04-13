

import { User } from "../api/models/user";
import { Caller } from "../api/api";

export const userService = {

  async get () : Promise<User | null> {
    try {
      const response = await Caller.get("/user");

      if (response.status !== 200) {
        return null;
      }

      return response.data as User;
    } catch (error) {
      return null
    }
  }

}