# companies/ — operativní záznam o firmách

Lehký CRM v markdownu. Jeden záznam na firmu, **napříč celým životním cyklem**
(prospekt → klient). Drží research, lidi, zápisy z meetingů a stav vztahu.

> Když se firma stane klientem, **technická realizace** problému jde do
> **vlastního repa** (viz `docs/vision.md` → Out of scope). Sem patří jen
> vztah: kdo tam je, co o nich víme, historie schůzek a stav.

## Struktura

```
companies/
  skoda-auto/
    skoda-auto.md      ← hlavní záznam (název = slug složky, kvůli fuzzy-find)
    (přílohy přidej, až budou: ares.json, *.pdf, screenshoty…)
```

- **Složka na firmu**, slug malými písmeny s pomlčkami (`skoda-auto`).
- **Hlavní soubor `<slug>.md`** — stejný název jako složka. Záměrně ne
  `README.md`: ve fuzzy-finderu chceš napsat „skoda" a trefit se, ne se
  prodírat padesáti stejně pojmenovanými soubory.
- **Přílohy** (PDF, surový výstup z `tools/ares`) do téže složky.

## Konvence záznamu

Frontmatter je strojově čitelný (pro agenta i grep), tělo je próza:

```markdown
---
name: Škoda Auto
ico: "00177041"
status: prospect        # prospect | contacted | meeting | won | lost
next_step: "poslat demo, do 2026-06-10"
---

## Research
## Lidi
## Zápisy z meetingů
```

## Pravidla

- **Stav žije ve frontmatteru (`status`), NE ve struktuře složek.** Žádné
  `prospects/` vs `clients/` — přesouvat soubory při každé změně stavu je
  otrava a rozbíjí odkazy. „Ukaž všechny prospekty" = grep/agent přes
  frontmatter.
- **`ico` je kanonické**, slug je jen pro lidi. `tools/ares` enrichuje podle IČO.
- **Nadčasová znalost sem nepatří** — ICP, playbook a SOP jsou ve `wiki/`.
