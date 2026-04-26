import { StateSchema } from "@langchain/langgraph";
import { z } from "zod";

const BomEntrySchema = z.object({
  cis_code: z.string(),
  name: z.string(),
  qty: z.number(),
  unit: z.string(),
});

export type BomEntry = z.infer<typeof BomEntrySchema>;

export const AgentState = new StateSchema({
  pdf_path: z.string(),
  drawing_number: z.string().nullable().default(null),
  bom: z.array(BomEntrySchema).nullable().default(null),
});
