import { z } from 'zod';

export const KeywordResearchSchema = z.object({
	query: z.string().optional(),
	country: z.string().min(2),
	period: z.string().min(1)
});

export type KeywordResearchInputs = z.infer<typeof KeywordResearchSchema>;
