import { z } from 'zod';

export const KeywordTrendExplorerSchema = z.object({
	query: z.string().optional(),
	country: z.string().min(2),
	language: z.string().optional().default('EN'),
	period: z.string().min(1)
});

export type KeywordTrendExplorerInputs = z.infer<
	typeof KeywordTrendExplorerSchema
>;
