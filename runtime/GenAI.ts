class _GenAI {
    private native: any;

    Generate(systemPrompt: string, prompt: string): string {
        let result = this.native.Generate(systemPrompt, prompt);

        return result
    }
}

declare const GenAI: _GenAI