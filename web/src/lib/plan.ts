import type { Blocker, Task, Tier, Today } from './api';
import { localDate } from './time';
import { available } from './execution';

export const tierLabels: Record<Tier, string> = {
	must: 'Must (Essencial)',
	should: 'Should (Importante)',
	could: 'Could (Se der)'
};

export interface Card {
	task: Task;
	rationale: string;
}

export type Tiers = Record<Tier, Card[]>;

export interface ActionCard {
	card: Card;
	tier?: Tier;
}

/** Another view of the current plan, not a separate queue to maintain. */
export function actionCards(today: Today, tasks: Task[], blockers: Blocker[], runningTaskID?: number): ActionCard[] {
	const grouped = tiers(today, tasks);
	const candidates: ActionCard[] = [
		...(['must', 'should', 'could'] as const).flatMap((tier) => grouped[tier].map((card) => ({card, tier}))),
		...unplanned(today, tasks).map((card) => ({card}))
	];
	// A live focus remains reachable even while the plan catches up with a recent change.
	const running = tasks.find((task) => task.id === runningTaskID && task.status === 'open');
	if (running && !candidates.some(({card}) => card.task.id === running.id)) candidates.unshift({card: {task: running, rationale: ''}});
	const seen = new Set<number>();
	const eligible = candidates.filter(({card}) => {
		if (seen.has(card.task.id) || (!available(card.task, tasks, blockers, today.date) && card.task.id !== running?.id)) return false;
		seen.add(card.task.id);
		return true;
	});
	const active = eligible.findIndex(({card}) => card.task.id === running?.id);
	if (active > 0) eligible.unshift(...eligible.splice(active, 1));
	return eligible;
}

/**
 * The plan's open tasks by tier, in plan order. A task the plan names but that is finished,
 * dropped, deleted, or pushed past today since is left out: it isn't something to do today any
 * more, and the plan only catches up at the next replan.
 */
export function tiers(today: Today, tasks: Task[]): Tiers {
	const byId = new Map(tasks.map((t) => [t.id, t]));
	const out: Tiers = { must: [], should: [], could: [] };
	for (const e of today.plan?.entries ?? []) {
		const task = byId.get(e.task_id);
		if (task?.status !== 'open' || (task.not_before ?? '') > today.date) continue;
		out[e.tier].push({ task, rationale: e.rationale });
	}
	return out;
}

/** Eligible tasks the plan leaves out, in the server's order. */
export function unplanned(today: Today, tasks: Task[]): Card[] {
	const byId = new Map(tasks.map((t) => [t.id, t]));
	return today.unplanned.flatMap((id) => {
		const task = byId.get(id);
		return task ? [{ task, rationale: '' }] : [];
	});
}

export interface Links {
	waitingOn: Task[];
	unblocks: Task[];
}

/**
 * What a task waits on and what waits on it. Only open tasks count: a finished blocker no longer
 * blocks, and a finished task no longer waits.
 */
export function links(id: number, blockers: Blocker[], tasks: Task[]): Links {
	const open = new Map(tasks.filter((t) => t.status === 'open').map((t) => [t.id, t]));
	const pick = (ids: number[]) => ids.flatMap((i) => open.get(i) ?? []);
	return {
		waitingOn: pick(blockers.filter((b) => b.task_id === id).map((b) => b.blocked_by_id)),
		unblocks: pick(blockers.filter((b) => b.blocked_by_id === id).map((b) => b.task_id))
	};
}

/** Tasks marked done on the given local date. */
export function doneOn(date: string, tasks: Task[]): Task[] {
	return tasks.filter((t) => t.status === 'done' && localDate(new Date(t.updated_at)) === date);
}

function count(n: number, one: string, many: string): string {
	return `${n} ${n === 1 ? one : many}`;
}

/** The Done today header: the day's win. */
export function winLine(done: number, meetings: number): string {
	if (done === 0 && meetings === 0) return 'Nothing checked off yet';
	const parts: string[] = [];
	if (done > 0) parts.push(count(done, 'thing done', 'things done'));
	if (meetings > 0) parts.push(count(meetings, 'meeting survived', 'meetings survived'));
	return `✓ ${parts.join(', ')}`;
}
