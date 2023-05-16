import { RgbaColor } from "react-colorful";
import rgbHex from "rgb-hex";
import hexRgb from "hex-rgb";

export const getRgbaCSS = (rgba: RgbaColor) => {
    const { r, g, b, a } = rgba;
    return `rgba(${r}, ${g}, ${b}, ${a})`;
}

export const getHexColor = (rgba: RgbaColor) => {
    const string = Array.from(Object.values(rgba)).join(",");

    const rgbaColor = "rgba(" + string + ")";
    return "#" + rgbHex(rgbaColor);
}

export const getRgbaObject = (hex: string): RgbaColor => {
    const {red: r,green: g, blue: b, alpha: a} = hexRgb(hex);
    return { r, g, b, a }
}