import { z } from 'zod';

export const SettingSchema = z.object({
	browserPath: z.string().optional(),
	browserVisible: z.boolean().default(false),
	sessionGoogle: z.string().optional(),
	sessionPinterest: z.string().min(1),
	proxyUrl: z.string().optional()
});

export type SettingInputs = z.infer<typeof SettingSchema>;
