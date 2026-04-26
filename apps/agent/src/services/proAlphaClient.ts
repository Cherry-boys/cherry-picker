const BASE = process.env.PROALPHA_URL!;
const TARGET = process.env.PROALPHA_TARGET_ID!;

type JobResult = {
  id: string;
  status: "completed" | "failed";
  result?: unknown;
  error?: string;
};

export async function submitJob(
  apiName: string,
  params: Record<string, string>
): Promise<string> {
  const res = await fetch(`${BASE}/api/targets/${TARGET}/jobs`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ api_name: apiName, params }),
  });
  if (!res.ok) throw new Error(`ProAlpha submitJob failed: ${res.status}`);
  const data = (await res.json()) as { id: string };
  return data.id;
}

export async function getJob(jobId: string): Promise<JobResult> {
  // The mock resolves synchronously; a real ProAlpha would need polling here.
  const res = await fetch(`${BASE}/api/targets/${TARGET}/jobs/${jobId}`);
  if (!res.ok) throw new Error(`ProAlpha getJob failed: ${res.status}`);
  return res.json() as Promise<JobResult>;
}
