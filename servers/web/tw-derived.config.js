module.exports = function ({ addBase, theme }) {
    const colors = theme("colors");

    const derived = {};

    Object.keys(colors).forEach((token) => {
        const base = colors[token];
        if (typeof base !== "string") return;

        derived[`:root`] = {
            [`--color-${token}-hover`]: `color-mix(in srgb, ${base} 90%, white)`,
            [`--color-${token}-dark`]: `color-mix(in srgb, ${base} 70%, black)`,
            [`--color-${token}-dark-hover`]: `color-mix(in srgb, ${base} 60%, black)`,
        };
    });

    addBase(derived);
};
