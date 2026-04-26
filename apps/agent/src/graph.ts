import { END, START, StateGraph } from "@langchain/langgraph";
import { AgentState } from "./state.js";
import { pdfExtractor } from "./nodes/pdfExtractor.js";
import { bomLookup } from "./nodes/bomLookup.js";

const builder = new StateGraph(AgentState)
  .addNode("pdf_extractor", pdfExtractor)
  .addNode("bom_lookup", bomLookup)
  .addEdge(START, "pdf_extractor")
  .addEdge("pdf_extractor", "bom_lookup")
  .addEdge("bom_lookup", END);

export const graph = builder.compile();
