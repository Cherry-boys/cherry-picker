import { prisma } from "../index.js";

const prospects = [
  {
    name: "Anna Nováková",
    company: "MedTech s.r.o.",
    email: "anna@medtech.example",
    phone: "+420 601 111 222",
    status: "new",
    notes: "Inbound from trade show.",
  },
  {
    name: "Jonas Berger",
    company: "Sterilex GmbH",
    email: "jonas.berger@sterilex.example",
    status: "contacted",
    notes: "Asked for a BOM automation demo.",
  },
];

async function main() {
  for (const data of prospects) {
    const existing = data.email
      ? await prisma.prospect.findFirst({ where: { email: data.email } })
      : null;
    if (existing) continue;
    await prisma.prospect.create({ data });
  }
  const count = await prisma.prospect.count();
  console.log(`Seed complete. Prospects in database: ${count}`);
}

main()
  .catch((err) => {
    console.error(err);
    process.exit(1);
  })
  .finally(() => prisma.$disconnect());
