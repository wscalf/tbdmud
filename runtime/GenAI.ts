declare class ChatMessage {
    Actor: string
    Body: string
}

class _GenAI {
    private native: any;

    Generate(systemPrompt: string, prompt: string): string {
        return this.native.Generate(systemPrompt, prompt);
    }
    Chat(systemPrompt: string, history: ChatMessage[], actor: string, pose: string): ChatMessage[] {
        return this.native.Chat(systemPrompt, history, actor, pose);
    }
}

declare const GenAI: _GenAI