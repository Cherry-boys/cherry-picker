#!/usr/bin/env bun

import { join } from "node:path";
import { Glob } from "bun";

const DOCS_DIR = join(import.meta.dir, "..", "docs");

const EXCLUDED_DIRS = new Set(["archive", "research", "plans", "internal"]);

function compactStrings(values: unknown[]): string[] {
	const result: string[] = [];
	for (const value of values) {
		if (value === null || value === undefined) {
			continue;
		}
		const normalized = String(value).trim();
		if (normalized.length > 0) {
			result.push(normalized);
		}
	}
	return result;
}

function isExcludedPath(relativePath: string): boolean {
	return relativePath.split(/[/\\]/).some((seg) => EXCLUDED_DIRS.has(seg));
}

function listMarkdownFiles(): string[] {
	const glob = new Glob("**/*.md");
	const files: string[] = [];
	for (const file of glob.scanSync({ cwd: DOCS_DIR })) {
		if (isExcludedPath(file)) {
			continue;
		}
		files.push(file);
	}
	return files.sort((a, b) => a.localeCompare(b));
}

async function extractMetadata(fullPath: string): Promise<{
	summary: string | null;
	readWhen: string[];
	error?: string;
}> {
	const content = await Bun.file(fullPath).text();

	if (!content.startsWith("---")) {
		return { summary: null, readWhen: [], error: "missing front matter" };
	}

	const endIndex = content.indexOf("\n---", 3);
	if (endIndex === -1) {
		return { summary: null, readWhen: [], error: "unterminated front matter" };
	}

	const frontMatter = content.slice(3, endIndex).trim();
	const lines = frontMatter.split("\n");

	let summaryLine: string | null = null;
	const readWhen: string[] = [];
	let collectingField: "read_when" | null = null;

	for (const rawLine of lines) {
		const line = rawLine.trim();

		if (line.startsWith("summary:")) {
			summaryLine = line;
			collectingField = null;
			continue;
		}

		if (line.startsWith("read_when:")) {
			collectingField = "read_when";
			const inline = line.slice("read_when:".length).trim();
			if (inline.startsWith("[") && inline.endsWith("]")) {
				try {
					const parsed = JSON.parse(inline.replace(/'/g, '"')) as unknown;
					if (Array.isArray(parsed)) {
						readWhen.push(...compactStrings(parsed));
					}
				} catch {
					// ignore malformed inline arrays
				}
			}
			continue;
		}

		if (collectingField === "read_when") {
			if (line.startsWith("- ")) {
				const hint = line.slice(2).trim();
				if (hint) {
					readWhen.push(hint);
				}
			} else if (line === "") {
			} else {
				collectingField = null;
			}
		}
	}

	if (!summaryLine) {
		return { summary: null, readWhen, error: "summary key missing" };
	}

	const summaryValue = summaryLine.slice("summary:".length).trim();
	const normalized = summaryValue
		.replace(/^['"]|['"]$/g, "")
		.replace(/\s+/g, " ")
		.trim();

	if (!normalized) {
		return { summary: null, readWhen, error: "summary is empty" };
	}

	return { summary: normalized, readWhen };
}

async function main(): Promise<void> {
	console.log("Listing all markdown files in docs folder:");

	const markdownFiles = listMarkdownFiles();

	for (const relativePath of markdownFiles) {
		const fullPath = join(DOCS_DIR, relativePath);
		const { summary, readWhen, error } = await extractMetadata(fullPath);
		if (summary) {
			console.log(`${relativePath} - ${summary}`);
			if (readWhen.length > 0) {
				console.log(`  Read when: ${readWhen.join("; ")}`);
			}
		} else {
			const reason = error ? ` - [${error}]` : "";
			console.log(`${relativePath}${reason}`);
		}
	}

	console.log(
		'\nReminder: keep docs up to date as behavior changes. When your task matches any "Read when" hint above (React hooks, cache directives, database work, tests, etc.), read that doc before coding, and suggest new coverage when it is missing.',
	);
}

await main();
