import type { Config } from 'tailwindcss';

export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],

	theme: {
		extend: {
			colors: {
				"c0": "#073605",
				"c1": "#0e4e0b",
				"c2": "#145b11",
				"c3": "#146c11",
				"c4": "#199515",
				"c5": "#000000"
			}
		}
	},

	plugins: []
} satisfies Config;
