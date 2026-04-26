import { Hono } from "hono";
import db from "../db.js";

type Material = {
  cis_code: string;
  name: string;
  unit: string;
  description: string | null;
};

type BomEntry = {
  cis_code: string;
  name: string;
  qty: number;
  unit: string;
};

type JobRecord = {
  id: string;
  status: "completed" | "failed";
  result?: unknown;
  error?: string;
};

const jobs = new Map<string, JobRecord>();

function getMaterial(params: Record<string, string>): Material {
  const lookup = params["belimed_nr"] ?? params["tyco_nr"] ?? "";
  const row = db
    .query<Material, [string, string]>(
      "SELECT cis_code, name, unit, description FROM materials WHERE belimed_nr = ? OR tyco_nr = ? LIMIT 1"
    )
    .get(lookup, lookup);
  if (!row) throw new Error(`Material not found for: ${lookup}`);
  return row;
}

function getBOM(params: Record<string, string>): BomEntry[] {
  const code = params["cis_product_code"] ?? "";
  return db
    .query<BomEntry, [string]>(
      `SELECT m.cis_code, m.name, be.qty, be.unit
       FROM bom_entries be
       JOIN materials m ON be.material_cis_code = m.cis_code
       WHERE be.cis_product_code = ?`
    )
    .all(code);
}

const router = new Hono();

router.post("/api/targets/:targetId/jobs", async (c) => {
  const body = await c.req.json<{
    api_name: string;
    params: Record<string, string>;
  }>();
  const id = crypto.randomUUID();

  let job: JobRecord;
  try {
    let result: unknown;
    if (body.api_name === "GetMaterial") {
      result = getMaterial(body.params);
    } else if (body.api_name === "GetBOM") {
      result = getBOM(body.params);
    } else {
      throw new Error(`Unknown api_name: ${body.api_name}`);
    }
    job = { id, status: "completed", result };
  } catch (err) {
    job = {
      id,
      status: "failed",
      error: err instanceof Error ? err.message : String(err),
    };
  }

  jobs.set(id, job);
  return c.json({ id, status: job.status });
});

router.get("/api/targets/:targetId/jobs/:jobId", (c) => {
  const jobId = c.req.param("jobId");
  const job = jobs.get(jobId);
  if (!job) return c.json({ error: "Job not found" }, 404);
  return c.json(job);
});

export default router;
