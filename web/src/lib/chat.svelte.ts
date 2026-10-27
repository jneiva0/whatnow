import { api, ApiError, type Message, type ToolUse } from './api';

/** The reply being written, built from deltas and tool chips until its stored message lands. */
export interface Draft {
	text: string;
	tools: ToolUse[];
}

export class Chat {
	messages = $state<Message[]>([]);
	draft = $state<Draft | null>(null);
	busy = $state(false);
	/** A gentle line under the input, such as when another tab is mid-reply. */
	notice = $state('');
	/** Reply batches undone, so their Undo buttons read "undone". */
	undone = $state<number[]>([]);
	open = $state(true);
	/** A reply landed while the dock was collapsed. */
	unread = $state(false);
	#loadGeneration = 0;
	#draftRevision = 0;
	// Request-local bookkeeping, not rendered/reactive collections.
	#duringLoad: Map<number, Message> | null = null;

	/** Recover history, retaining confirmed messages that arrive while it is in flight. */
	async load(): Promise<void> {
		const generation = ++this.#loadGeneration;
		const draftRevision = this.#draftRevision;
		const received = new Map<number, Message>();
		this.#duringLoad = received;
		try {
			const history = await api.chat();
			if (generation !== this.#loadGeneration) return;
			const merged = new Map(history.map((m) => [m.id, m]));
			for (const [id, message] of received) merged.set(id, message);
			this.messages = [...merged.values()].sort((a, b) => a.id - b.id);
			// Clear a stranded draft, but never a new reply arriving over the recovered stream.
			if (draftRevision === this.#draftRevision) this.draft = null;
		} finally {
			if (generation === this.#loadGeneration) this.#duringLoad = null;
		}
	}

	delta(text: string): void {
		++this.#draftRevision;
		this.draft ??= { text: '', tools: [] };
		this.draft.text += text;
	}

	tool(use: ToolUse): void {
		++this.#draftRevision;
		this.draft ??= { text: '', tools: [] };
		this.draft.tools.push(use);
	}

	/** A stored message is authoritative: it replaces the draft it finishes. */
	message(m: Message): void {
		this.#duringLoad?.set(m.id, m);
		const i = this.messages.findIndex((x) => x.id === m.id);
		if (i >= 0) this.messages[i] = m;
		else this.messages.push(m);
		// SSE and the POST response carry the same persisted message. Re-delivery
		// must not finish a newer draft or mark the old reply unread again.
		if (m.role === 'assistant' && i < 0) {
			++this.#draftRevision;
			this.draft = null;
			if (!this.open) this.unread = true;
		}
	}

	markUndone(batch: number): void {
		this.undone.push(batch);
	}

	toggle(): void {
		this.open = !this.open;
		if (this.open) this.unread = false;
	}

	/** Sends a message. Returns false when it wasn't taken, so the input keeps the text. */
	async send(text: string): Promise<boolean> {
		this.busy = true;
		this.notice = '';
		try {
			this.message(await api.send(text));
			return true;
		} catch (err) {
			if (err instanceof ApiError && err.status === 409) {
				this.notice = 'One sec, still answering.';
			} else {
				this.notice = err instanceof Error ? err.message : String(err);
			}
			return false;
		} finally {
			this.busy = false;
		}
	}
}

export const chat = new Chat();
