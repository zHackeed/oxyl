import { Logger } from "tslog";

export const logger = new Logger({
  name: "oxyl-websocket",
  type: "pretty",
});