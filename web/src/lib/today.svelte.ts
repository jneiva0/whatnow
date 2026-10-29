import { api, type Blocker, type CalendarEvent, type Task, type Today } from './api';
import { toasts } from './toasts.svelte';

const debounce = 100;

/** Today's state, refetched whole whenever the server says something changed. */
export class TodayState {
	today = $state<Today | null>(null);
	events = $state<CalendarEvent[]>([]);
	tasks = $state<Task[]>([]);
	blockers = $state<Blocker[]>([]);
	revision = $state(0);
	#timer: ReturnType<typeof setTimeout> | undefined;
	#generation = 0;

	async refresh(): Promise<void> {
		const generation = ++this.#generation;
		const today = await api.today();
		if (generation !== this.#generation) return;
		const [events, tasks, blockers] = await Promise.all([
			api.events(today.date),
			api.tasks(),
			api.blockers()
		]);
		// A slower earlier read must not replace a newer stream-triggered snapshot.
		if (generation !== this.#generation) return;
		this.today = today;
		this.events = events;
		this.tasks = tasks;
		this.blockers = blockers;
		this.revision++;
	}

	/** Refreshes once a burst of changes settles. */
	changed(): void {
		clearTimeout(this.#timer);
		this.#timer = setTimeout(() => this.refresh().catch((err) => toasts.error(err)), debounce);
	}
}

export const today = new TodayState();
