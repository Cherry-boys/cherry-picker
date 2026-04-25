import { END, START, StateGraph, Annotation } from "@langchain/langgraph";

const StateAnnotation = Annotation.Root({
  message: Annotation<string>,
});

const greet = async (state: typeof StateAnnotation.State) => {
  return { message: `Hello from the graph! Input was: ${state.message}` };
};

const builder = new StateGraph(StateAnnotation)
  .addNode("greet", greet)
  .addEdge(START, "greet")
  .addEdge("greet", END);

export const graph = builder.compile();
