import { Hono } from "hono";
import "./db.js";
import jobsRouter from "./routes/jobs.js";

const app = new Hono();
app.route("/", jobsRouter);

const port = 3001;
console.log(`proalpha-mock listening on :${port}`);

export default {
  port,
  fetch: app.fetch,
};
