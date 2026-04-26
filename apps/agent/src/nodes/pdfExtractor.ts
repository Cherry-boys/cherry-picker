import { readFile } from "fs/promises";
import { HumanMessage } from "@langchain/core/messages";
import { AgentState } from "../state.js";
import { mainModel } from "../models/index.js";

export async function pdfExtractor(
  state: typeof AgentState.State
): Promise<Partial<typeof AgentState.State>> {
  const bytes = await readFile(state.pdf_path);
  const base64 = bytes.toString("base64");

  const response = await mainModel.invoke([
    new HumanMessage({
      content: [
        {
          type: "document",
          source: {
            type: "base64",
            media_type: "application/pdf",
            data: base64,
          },
        } as never,
        {
          type: "text",
          text: "Extract the drawing number from this technical drawing. Return only the number, nothing else.",
        },
      ],
    }),
  ]);

  const drawing_number = (response.content as string).trim();
  return { drawing_number };
}
