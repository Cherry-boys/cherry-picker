import { ChatAnthropic } from "@langchain/anthropic";

export const mainModel = new ChatAnthropic({
  model: "claude-sonnet-4-6",
});
