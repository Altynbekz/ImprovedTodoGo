# AGENTS.md — Frontend Architecture & Design System

## 1. Visual Identity & Aesthetic Direction
* **Concept:** Brutalist Swiss Editorial meets Generative Kinetic Grid ("A modern yin and yang. Two values, nothing else.").
* **Strict Dichotomy:** Every visual element on the screen must strictly adhere to a binary state (`0` or `1` / `True` or `False`). No gradients, no drop shadows, no intermediate grays, and no rounded corners (`border-radius: 0px` globally).
* **Color Palette:**
  * `--color-primary: #0014FF` (Electric International Klein Blue)
  * `--color-surface: #FFFFFF` (Stark Pure White)
  * *Rule:* Elements must strictly invert between `--color-primary` and `--color-surface`.

---

## 2. Typography System
* **Display / Hero Font:** Heavy Neo-Grotesque Sans-Serif (`Inter` 800/900, `Helvetica Now Display`, or `Archivo Black`).
  * Tight tracking (`letter-spacing: -0.04em`), tight line-height (`0.95`), high visual mass.
* **UI / Technical / Micro-copy Font:** Crisp Monospace (`JetBrains Mono`, `Space Mono`, or `IBM Plex Mono`).
  * Used for metadata, grid symbols, buttons, and technical captions.
  * Uppercase labels use wide tracking (`letter-spacing: 0.12em`, `font-size: 10px`–`11px`).

---

## 3. Core Architecture: The Interactive Boolean Grid Canvas
The entire viewport (`100vw` x `100vh`) acts as a living, full-bleed 2D grid matrix rendered via HTML5 `<canvas>` (for 60fps performance) overlaid with crisp, grid-snapped DOM UI blocks.

### Grid Mechanics
* **Cell Dimensions:** Uniform square/rectangular grid cells (e.g., `24px x 24px` or `28px x 28px`), dynamically calculated on window resize to fill the viewport cleanly.
* **Binary States:**
  * **State `0` (False):** Background `#FFFFFF`, Foreground symbol `#0014FF`.
  * **State `1` (True):** Background `#0014FF`, Foreground symbol `#FFFFFF`.
* **Symbol Modes (Controlled by Top-Right Toggle):**
  1. `+/-` Mode: State `0` renders `-`, State `1` renders `+`.
  2. `0/1` Mode: State `0` renders `0`, State `1` renders `1`.
  3. `>/<` Mode (Default): State `0` renders `<`, State `1` renders `>`.
* **Organic Noise & Wave Physics:**
  * Base state is driven by slow, organic 2D Simplex/Perlin noise creating "continents" of `1`s and `0`s.
  * **Click Interaction ("CLICK THE FIELD TO SEND A WAVE THROUGH IT"):** Clicking anywhere on the canvas triggers an expanding radial shockwave that temporarily inverts the boolean state of every cell it passes through.
  * **Play/Pause State:** Respects the top-right `pause` / `play` toggle, freezing or resuming the noise evolution and wave propagation.

---

## 4. Layout & Grid-Snapped UI Components
All UI overlays must feel carved out of the underlying grid, using solid `#FFFFFF` or `#0014FF` rectangular blocks with sharp `90deg` edges and zero transparency.

### A. Top Navigation Bar (`position: fixed; top: 24px; left: 24px; right: 24px`)
* **Top-Left Brand Badge:**
  * Solid `#FFFFFF` rectangle (`padding: 8px 14px`).
  * Text: `boolean` (Regular/Medium) + `Sequence` (Bold) in `#0014FF`.
  * Sub-caption directly below on transparent/grid background: `NICEUNDERGROUND` `GENERATIVE VISUALS` (Monospace, uppercase, 10px, tracking `0.15em`, adapting color to contrast with underlying cells or enclosed in tight micro-pills).
* **Top-Right Control Cluster:**
  * A contiguous horizontal segmented bar (`border: 1px solid #0014FF` or high-contrast solid blocks):
    * Mode Buttons: `[ +/- ]`, `[ 0/1 ]`, `[ >/< ]`. Active button has `#0014FF` fill and `#FFFFFF` text; inactive buttons have `#FFFFFF` fill and `#0014FF` text.
    * Action Button (separated by a `12px` gap): `[ || pause ]` / `[ ▶ play ]` in a solid `#FFFFFF` block with `#0014FF` text.

### B. Right-Side Editorial Callout
* Positioned in the upper-right quadrant (`top: 22%`, `right: 24px`).
* **Card:** Solid `#FFFFFF` rectangular container (`padding: 16px 20px`, `max-width: 280px`).
* **Copy:** Clean sans-serif (`15px`, `#0014FF`): *"A modern yin and yang. Two values, nothing else."*
* **Micro-instruction below card:** Monospace uppercase (`9px`–`10px`): `> CLICK THE FIELD TO SEND A WAVE THROUGH IT`

### C. Center-Left Hero Typography ("Cutout" Blocks)
* Positioned in the lower-middle left area, spanning ~70–85% of viewport width on desktop.
* **Staggered Block Treatment:** Instead of floating text over the grid, each line sits inside a tight, solid `#FFFFFF` bounding box that masks the grid behind it, creating a stepped, brutalist silhouette:
  * Line 1 Box: `boolean` (Flush left, massive heavy sans-serif, `#0014FF`).
  * Line 2 Box: `Sequence` (Indented slightly right relative to Line 1, massive heavy sans-serif, `#0014FF`).
* *Optional Advanced Polish:* Allow grid cells adjacent to the hero boxes to dynamically react or "wrap" around the stepped edges.

### D. Vertical Edge Indicator (Right Margin)
* Fixed vertically centered on the far right edge (`right: 16px`).
* Small vertical stack of micro-squares (`6px x 6px`, `1px solid #0014FF`) indicating scroll/section progress (top square filled `#0014FF`, remaining hollow `#FFFFFF`), paired with rotated vertical micro-text (`90deg`).

### E. Bottom Bar (`position: fixed; bottom: 24px; left: 24px; right: 24px`)
* **Bottom-Left Scroll Prompt:**
  * Solid `#FFFFFF` square/rectangular badge (`padding: 8px 12px`), monospace lowercase `scroll` in `#0014FF`.
* **Bottom-Right Poetic Footer:**
  * Solid `#FFFFFF` horizontal strip (`padding: 6px 12px`).
  * Copy: `contrasts that swap places,` `[ > ]` `never merge` `[ < ]` where `[ > ]` and `[ < ]` are rendered as inverted `#0014FF` square badges with `#FFFFFF` glyphs.

---

## 5. Coding Standards & Implementation Rules
1. **Zero-Radius Rule:** Never use `rounded-*` classes or `border-radius`. Everything is strictly orthogonal.
2. **Performance First:** Render the background grid on a single `<canvas>` element using `requestAnimationFrame`. Do not create thousands of individual DOM `<div>` nodes for the grid cells.
3. **High-DPI Crispness:** Scale the `<canvas>` by `window.devicePixelRatio` so monospace symbols (`>`, `<`, `0`, `1`, `+`, `-`) render razor-sharp on Retina/4K displays.
4. **Responsive Adaptation:**
   * On mobile/tablet viewports (`< 768px`), scale down the Hero `boolean` / `Sequence` font size using `clamp(3.2rem, 13vw, 9rem)` while preserving the staggered two-box cutout layout.
   * Keep the top-right control cluster compact and touch-friendly (`min-height: 36px`).
5. **Accessibility:** Ensure keyboard focus states use high-contrast `2px solid` outlines in the opposing binary color, and honor `prefers-reduced-motion` by starting the grid in the `paused` state if requested by the OS.