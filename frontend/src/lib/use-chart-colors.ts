"use client";

import { useTheme } from "next-themes";
import { useEffect, useState } from "react";

const LIGHT_DEFAULTS = {
  stroke: "#18181b",
  fill: "#3f3f46",
  bar: "#52525b",
  barAlt: "#71717a",
  tick: "#71717a",
  gradient: "#18181b",
  tooltipBg: "#ffffff",
  tooltipBorder: "#e4e4e7",
  tooltipText: "#18181b",
};

const DARK_DEFAULTS = {
  stroke: "#d4d4d8",
  fill: "#a1a1aa",
  bar: "#a1a1aa",
  barAlt: "#d4d4d8",
  tick: "#a1a1aa",
  gradient: "#d4d4d8",
  tooltipBg: "#18181b",
  tooltipBorder: "#3f3f46",
  tooltipText: "#fafafa",
};

function readCSSVars(): typeof LIGHT_DEFAULTS {
  if (typeof window === "undefined") return LIGHT_DEFAULTS;
  const style = getComputedStyle(document.documentElement);
  return {
    stroke: style.getPropertyValue("--chart-stroke").trim() || LIGHT_DEFAULTS.stroke,
    fill: style.getPropertyValue("--chart-fill").trim() || LIGHT_DEFAULTS.fill,
    bar: style.getPropertyValue("--chart-bar").trim() || LIGHT_DEFAULTS.bar,
    barAlt: style.getPropertyValue("--chart-bar-alt").trim() || LIGHT_DEFAULTS.barAlt,
    tick: style.getPropertyValue("--chart-tick").trim() || LIGHT_DEFAULTS.tick,
    gradient: style.getPropertyValue("--chart-gradient").trim() || LIGHT_DEFAULTS.gradient,
    tooltipBg: style.getPropertyValue("--tooltip-bg").trim() || LIGHT_DEFAULTS.tooltipBg,
    tooltipBorder: style.getPropertyValue("--tooltip-border").trim() || LIGHT_DEFAULTS.tooltipBorder,
    tooltipText: style.getPropertyValue("--tooltip-text").trim() || LIGHT_DEFAULTS.tooltipText,
  };
}

export function useChartColors() {
  const { resolvedTheme } = useTheme();
  const [colors, setColors] = useState(LIGHT_DEFAULTS);

  useEffect(() => {
    // Small delay to let the CSS class update on <html>
    const id = requestAnimationFrame(() => setColors(readCSSVars()));
    return () => cancelAnimationFrame(id);
  }, [resolvedTheme]);

  return colors;
}
