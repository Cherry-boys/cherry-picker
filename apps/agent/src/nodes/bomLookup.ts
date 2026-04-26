import { AgentState, BomEntry } from "../state.js";
import { submitJob, getJob } from "../services/proAlphaClient.js";

export async function bomLookup(
  state: typeof AgentState.State
): Promise<Partial<typeof AgentState.State>> {
  const drawing_number = state.drawing_number;
  if (!drawing_number) throw new Error("drawing_number is missing in state");

  const jobId = await submitJob("GetBOM", { cis_product_code: drawing_number });
  const job = await getJob(jobId);

  if (job.status !== "completed") {
    throw new Error(`ProAlpha GetBOM failed: ${job.error ?? "unknown error"}`);
  }

  const bom = job.result as BomEntry[];
  return { bom };
}
