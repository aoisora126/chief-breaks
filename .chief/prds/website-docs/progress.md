## Codebase Patterns
- VitePress docs live in `/docs` directory with `.vitepress/config.ts` for configuration
- Use `npm run docs:dev` to start dev server, `npm run docs:build` for production build
- Base URL is `/chief/` for GitHub Pages project site
- Custom theme lives in `docs/.vitepress/theme/` with `index.ts` extending DefaultTheme
- Vue components for theme go in `docs/.vitepress/theme/components/` directory
- Custom layouts extend DefaultTheme.Layout and can use slots like `#home-hero-before`
- Links in Vue components must include base path (e.g., `/chief/guide/` not `/guide/`)
- Tailwind CSS v4 uses `@import "tailwindcss"` directive (not v3's `@tailwind` directives)
- Tailwind v4 plugin is configured via VitePress's `vite.plugins` option in config.ts
- Tokyo Night color palette defined in `tailwind.css` with both Tailwind v4 `@theme` and VitePress CSS variables
- Site is forced dark mode via `appearance: 'force-dark'` in config.ts
- Code blocks use Shiki's `tokyo-night` theme via `markdown.theme` in config.ts
- Links in markdown files should NOT include base path (use `/guide/quick-start` not `/chief/guide/quick-start`)
- Links in Vue components MUST include base path (use `/chief/guide/quick-start` not `/guide/quick-start`)

---

## 2026-01-28 - US-001
- **What was implemented**: VitePress project scaffolded in /docs directory
- **Files changed**:
  - `.gitignore` - added node_modules and VitePress build artifacts
  - `docs/package.json` - created with dev/build/preview scripts
  - `docs/package-lock.json` - dependency lockfile
  - `docs/.vitepress/config.ts` - site configuration with title "Chief", tagline "Autonomous PRD Agent", base URL `/chief/`
  - `docs/index.md` - landing page with hero layout
  - `docs/guide/index.md` - placeholder getting started page
- **Learnings for future iterations:**
  - VitePress v1.6.x uses Vite under the hood, configuration is in `.vitepress/config.ts`
  - The `base` option in config must match GitHub Pages project path (e.g., `/chief/`)
  - VitePress creates `.vitepress/cache/` and `.vitepress/dist/` directories that should be gitignored
  - Dev server runs on port 5173 by default
---

## 2026-01-28 - US-002
- **What was implemented**: Tailwind CSS v4 integration with VitePress
- **Files changed**:
  - `docs/package.json` - added tailwindcss and @tailwindcss/vite dependencies
  - `docs/package-lock.json` - updated lockfile with new dependencies
  - `docs/.vitepress/config.ts` - added Tailwind v4 Vite plugin via `vite.plugins` option
  - `docs/.vitepress/theme/index.ts` - custom theme extending DefaultTheme and importing tailwind.css
  - `docs/.vitepress/theme/tailwind.css` - CSS file with `@import "tailwindcss"` directive
- **Learnings for future iterations:**
  - Tailwind CSS v4 uses `@import "tailwindcss"` directive instead of v3's `@tailwind base/components/utilities`
  - VitePress custom themes go in `.vitepress/theme/` with `index.ts` as entry point
  - The theme must re-export DefaultTheme to preserve VitePress default styling
  - Vite plugins are added to VitePress via the `vite` config option in `.vitepress/config.ts`
  - Tailwind v4 is purely CSS-based with no separate config file needed for basic setup
---

## 2026-01-28 - US-003
- **What was implemented**: Tokyo Night dark theme for the documentation site
- **Files changed**:
  - `docs/.vitepress/config.ts` - added `appearance: 'force-dark'` and `markdown.theme: 'tokyo-night'` for code blocks
  - `docs/.vitepress/theme/tailwind.css` - extensive Tokyo Night color palette and VitePress CSS variable overrides
- **Learnings for future iterations:**
  - VitePress uses `appearance: 'force-dark'` to force dark mode and hide the theme toggle
  - Shiki (VitePress's syntax highlighter) has built-in `tokyo-night` theme - just set `markdown.theme: 'tokyo-night'`
  - VitePress CSS variables are organized by category: `--vp-c-brand-*`, `--vp-c-bg-*`, `--vp-c-text-*`, etc.
  - Tailwind v4 uses `@theme` directive to define custom color utilities (e.g., `--color-tokyo-bg` becomes `bg-tokyo-bg`)
  - VitePress class names like `.VPSidebar`, `.VPNav`, `.VPContent` can be styled directly
  - Custom block (tip, warning, danger) colors use `--vp-c-tip-*`, `--vp-c-warning-*`, `--vp-c-danger-*` variables
  - Force dark mode for non-.dark html with duplicate CSS variables to prevent flash of light theme
---

## 2026-01-28 - US-004
- **What was implemented**: Landing page hero section with animated terminal
- **Files changed**:
  - `docs/.vitepress/theme/components/Hero.vue` - new custom Hero component with headline, terminal animation, install command, and CTA buttons
  - `docs/.vitepress/theme/HomeLayout.vue` - custom layout that hides default VitePress hero and uses custom Hero component
  - `docs/.vitepress/theme/index.ts` - updated to use HomeLayout as the main layout
- **Learnings for future iterations:**
  - VitePress allows custom layouts via named slots like `#home-hero-before`, `#home-hero-info`, etc.
  - To completely replace the default hero, use a custom Layout component that extends DefaultTheme.Layout
  - Hide default VitePress hero with `.VPHome .VPHero { display: none !important; }`
  - CSS animations with `animation-delay` can create sequenced typing/fadeIn effects for terminal output
  - Vue components in VitePress theme go in `.vitepress/theme/components/` directory
  - Links in Vue components should use the full base path (e.g., `/chief/guide/` not `/guide/`)
---

## 2026-01-28 - US-005
- **What was implemented**: Landing page "How It Works" section with three-step visual workflow
- **Files changed**:
  - `docs/.vitepress/theme/components/HowItWorks.vue` - new component with three steps: Write PRD → Chief Runs Loop → Code Gets Built
  - `docs/.vitepress/theme/HomeLayout.vue` - updated to include HowItWorks component after Hero
- **Learnings for future iterations:**
  - Landing page sections are added to HomeLayout via the `#home-hero-before` slot after other components
  - Tokyo Night color variables can be used directly in component styles (e.g., `#7aa2f7` for accent, `#bb9af7` for purple, `#9ece6a` for green)
  - SVG icons from Feather Icons work well for step illustrations
  - Flexbox with `flex-direction: column` on mobile and row on desktop handles responsive step layouts
  - Step connectors (arrows) should rotate 90 degrees on mobile to maintain visual flow
---

## 2026-01-28 - US-006
- **What was implemented**: Landing page "Key Features" section with four feature cards in a grid layout
- **Files changed**:
  - `docs/.vitepress/theme/components/Features.vue` - new component with 4 feature cards: Single Binary, Self-Contained State, Works Anywhere, Beautiful TUI
  - `docs/.vitepress/theme/HomeLayout.vue` - updated to include Features component after HowItWorks
- **Learnings for future iterations:**
  - CSS Grid with `grid-template-columns: repeat(2, 1fr)` creates a 2-column layout that gracefully collapses to 1 column on mobile
  - Different hover border colors for each card can be achieved with `:nth-child(n)` selectors
  - Use `rgba()` for semi-transparent background colors on feature icons (e.g., `rgba(122, 162, 247, 0.1)`)
  - Alternate section backgrounds between `#1a1b26` and `#16161e` for visual separation
  - The cyan color for Tokyo Night is `#7dcfff` (useful for UI/TUI related icons)
---

## 2026-01-28 - US-007
- **What was implemented**: Landing page footer with CTA section
- **Files changed**:
  - `docs/.vitepress/theme/components/Footer.vue` - new footer component with "Ready to automate your PRDs?" CTA, links to quick start guide and GitHub, and copyright notice
  - `docs/.vitepress/theme/HomeLayout.vue` - updated to include Footer component after Features
- **Learnings for future iterations:**
  - Footer sections can be placed in the `#home-hero-before` slot along with other landing page sections
  - Use `border-top` to visually separate footer from content above
  - CTA buttons follow same styling pattern as Hero: primary (filled) and secondary (outlined)
  - Dynamic year in copyright: `{{ new Date().getFullYear() }}` works in Vue template
  - Footer uses darker `#16161e` background to contrast with `#1a1b26` features section above
---

## 2026-01-28 - US-008
- **What was implemented**: Navigation and sidebar structure with all documentation pages
- **Files changed**:
  - `docs/.vitepress/config.ts` - added top nav (Home, Docs, GitHub link) and full sidebar configuration
  - `docs/guide/index.md` - updated getting started landing page with links to subpages
  - `docs/guide/quick-start.md` - new quick start guide
  - `docs/guide/installation.md` - new detailed installation guide
  - `docs/concepts/how-it-works.md` - new overview of how Chief works
  - `docs/concepts/ralph-loop.md` - new deep dive into the Ralph Loop
  - `docs/concepts/prd-format.md` - new PRD format documentation
  - `docs/concepts/chief-directory.md` - new .chief directory guide
  - `docs/reference/cli.md` - new CLI reference
  - `docs/reference/configuration.md` - new configuration docs
  - `docs/reference/prd-schema.md` - new PRD schema reference
  - `docs/troubleshooting/common-issues.md` - new common issues guide
  - `docs/troubleshooting/faq.md` - new FAQ page
- **Learnings for future iterations:**
  - VitePress sidebar is configured via `themeConfig.sidebar` array with `text` and `items` for sections
  - Links in markdown files should NOT include base path (VitePress handles it automatically)
  - VitePress validates dead links during build - useful for catching broken internal links
  - Navigation items can be simple links or have children for dropdowns
  - Mobile navigation automatically collapses to hamburger menu (built into VitePress)
---

## 2026-01-28 - US-009
- **What was implemented**: Comprehensive Quick Start guide with all installation options and step-by-step instructions
- **Files changed**:
  - `docs/guide/quick-start.md` - expanded from placeholder to full guide with prerequisites, installation options (Homebrew, install script, from source), step-by-step workflow, TUI explanation with keyboard controls, and next steps links
- **Learnings for future iterations:**
  - VitePress `::: code-group` syntax creates tabbed code blocks for showing multiple installation options
  - Tables in markdown work well for keyboard shortcuts reference (pipe-separated columns with header row)
  - VitePress custom blocks `::: tip`, `::: info`, `::: warning` are useful for highlighting prerequisites and notes
  - Quick start guides should focus on "get running fast" with links to deeper docs, not exhaustive detail
---

## 2026-01-28 - US-010
- **What was implemented**: Detailed installation guide with all platform coverage
- **Files changed**:
  - `docs/guide/installation.md` - expanded from basic guide to comprehensive installation reference with prerequisites at top, Homebrew with update instructions, install script options table, complete platform matrix, platform-specific code tabs for manual download, building from source with version embedding, and thorough verification section
- **Learnings for future iterations:**
  - `::: code-group` is excellent for platform-specific installation commands (one tab per platform)
  - Architecture detection commands: `uname -m` returns `arm64`/`x86_64` on macOS, `aarch64`/`x86_64` on Linux
  - VitePress `::: info` blocks are good for notes about PATH configuration
  - `::: warning` blocks work well for troubleshooting tips at the end of installation sections
  - Tables are good for option flags documentation (Option | Description | Example format)
---

## 2026-01-28 - US-012
- **What was implemented**: Comprehensive Ralph Loop deep dive page with detailed step-by-step explanation
- **Files changed**:
  - `docs/concepts/ralph-loop.md` - expanded from basic overview to full deep dive with:
    - Updated blog post link to actual URL (larswadefalk.com)
    - Enhanced Mermaid flowchart with story selection, iteration limits, and Tokyo Night color styling
    - 7 detailed steps (Read State, Select Next Story, Build Prompt, Invoke Claude Code, Stream & Parse Output, Watch for Completion Signal, Update and Continue)
    - Tables showing files read and what Chief learns from each
    - Story selection logic explanation (priority sorting, inProgress handling)
    - Simplified example of the embedded prompt Claude receives
    - ASCII diagram showing stream-json output format with message types
    - Detailed explanation of `<chief-complete/>` signal and what it implies
    - Iteration limits section with scenario table and troubleshooting tips
    - "What's Next" links to related docs
- **Learnings for future iterations:**
  - Mermaid flowcharts support `style` directives for Tokyo Night colors (fill, stroke, color)
  - Use `([text])` for stadium-shaped (rounded) nodes in Mermaid for start/end states
  - Tables are effective for showing file-to-purpose mappings
  - Code blocks with ASCII box drawing characters create effective stream visualizations
  - Numbered lists with bold step names and sub-bullets create scannable deep dive content
---

## 2026-01-28 - US-011
- **What was implemented**: Enhanced "How Chief Works" overview page with comprehensive documentation
- **Files changed**:
  - `docs/concepts/how-it-works.md` - expanded from basic placeholder to full overview with:
    - High-level explanation of autonomous agent concept vs traditional interactive prompting
    - Improved ASCII diagram showing: User → PRD → Chief → Claude → Code pipeline
    - Component table explaining each part of the system
    - Detailed 7-step iteration flow explaining how stories are processed
    - New "Conventional Commits" section showing commit message format
    - New "Progress Tracking" section explaining progress.md and learnings
    - Link to blog post in a tip callout at the top
    - Updated links to related docs using em-dash formatting
- **Learnings for future iterations:**
  - VitePress `::: tip` blocks with custom headers work well for important links/callouts
  - ASCII art diagrams should fit within 80 characters for readability
  - Tables are effective for explaining system components with Role descriptions
  - Numbered lists with bold step names create scannable process documentation
  - Using em-dashes (—) for link descriptions creates consistent visual style
---

## 2026-01-28 - US-013
- **What was implemented**: Comprehensive .chief directory guide with detailed structure and file explanations
- **Files changed**:
  - `docs/concepts/chief-directory.md` - expanded from basic placeholder to full guide with:
    - Enhanced directory tree showing project context (not just `.chief/` in isolation)
    - Detailed `prds/` subdirectory explanation with CLI usage
    - Expanded file explanations with field tables for `prd.json`, example entries for `progress.md`, and usage context for each file
    - "Self-Contained by Design" section emphasizing no global config, no conflicts, no cleanup
    - Enhanced portability section with multiple examples (move, clone, remote)
    - Multiple PRDs section with practical examples
    - Git considerations with tables for commit/ignore decisions and `.gitignore` pattern
    - "What's Next" navigation links
- **Learnings for future iterations:**
  - VitePress doesn't have built-in `gitignore` language for syntax highlighting — it falls back to `txt` (cosmetic warning only, not a build error)
  - Tables with Yes/No columns and "Why" explanations are effective for commit/ignore decisions
  - Directory trees that show surrounding project context (e.g., `src/`, `package.json`) help users understand where `.chief/` fits
  - `::: tip` blocks work well for collaborative workflow notes
---
