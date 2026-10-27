import { describe, expect, it } from 'vitest';
import type { FocusSession, Step, Task } from './api';
import { elapsedMs, orderedSteps, timerText, available, actionClock, actionActual, taskActual, taskEstimate, focusTarget, reorderStepIDs } from './execution';

const session: FocusSession = {id: 1, task_id: 2, step_id: null, state: 'running', segment_started_at: '2026-10-25T10:00:00Z', accumulated_ms: 120000, created_at: '2026-10-25T09:58:00Z', ended_at: null};

describe('execution clock', () => {
 it('derives elapsed from server time across a suspended tab and keeps counting past zero', () => {
  expect(elapsedMs(session, Date.parse('2026-10-25T10:10:00Z'))).toBe(720000);
  expect(actionClock(5, 720000)).toEqual({text: '+07:00', label: 'Tempo excedido'});
  expect(session.state).toBe('running');
 });
 it('holds accumulated time when paused or ended and ignores clock reversal', () => {
  expect(elapsedMs({...session, state: 'paused'}, Date.parse('2026-10-25T11:00:00Z'))).toBe(120000);
  expect(elapsedMs({...session, state: 'ended'}, Date.parse('2026-10-25T11:00:00Z'))).toBe(120000);
  expect(elapsedMs(session, Date.parse('2026-10-25T09:00:00Z'))).toBe(120000);
 });
 it('formats long actual sessions without wrapping minutes to zero', () => {
  expect(timerText(3661000)).toBe('61:01');
  expect(timerText(-1000)).toBe('00:00');
 });
});

describe('execution targets', () => {
 it('leaves deleted historical steps out and orders pending and completed siblings together', () => {
  const step = {task_id: 2, description: 'Action', estimate_minutes: null, done: false, deleted_at: null};
  const steps: Step[] = [{...step, id: 1, position: 2}, {...step, id: 2, position: 0, done: true}, {...step, id: 3, position: 1, deleted_at: '2026-10-25'}, {...step, id: 4, task_id: 8, position: 0}];
  expect(orderedSteps(steps, 2).map((s) => s.id)).toEqual([2, 1]);
 });
 it('keeps blockers, deferred tasks and closed tasks out of the manual chooser', () => {
  const task = {id: 1, status: 'open', not_before: null} as Task;
  const blocker = {id: 2, status: 'open'} as Task;
  expect(available(task, [task, blocker], [{task_id: 1, blocked_by_id: 2}], '2026-10-25')).toBe(false);
  expect(available(task, [task, {...blocker, status: 'done'}], [{task_id: 1, blocked_by_id: 2}], '2026-10-25')).toBe(true);
  expect(available({...task, not_before: '2026-10-26'}, [task], [], '2026-10-25')).toBe(false);
  expect(available({...task, status: 'done'}, [task], [], '2026-10-25')).toBe(false);
 });
});


describe('action estimate and total time', () => {
 it('counts down the action estimate, shows excess, and uses elapsed with no estimate', () => {
  expect(actionClock(2, 110000)).toEqual({text: '00:10', label: 'Tempo restante'});
  expect(actionClock(2, 130000)).toEqual({text: '+00:10', label: 'Tempo excedido'});
  expect(actionClock(null, 130000)).toEqual({text: '02:10', label: 'Tempo dedicado'});
 });
 it('sums repeated action sessions once and parent includes direct plus children without duplication', () => {
  const sessions = [{...session, id: 1, step_id: null, state: 'ended' as const, accumulated_ms: 20000}, {...session, id: 2, step_id: 9, state: 'ended' as const, accumulated_ms: 30000}, {...session, id: 3, step_id: 9, state: 'ended' as const, accumulated_ms: 40000}];
  expect(actionActual(sessions, 2, 9, Date.now())).toBe(70000);
  expect(taskActual(sessions, 2, Date.now())).toBe(90000);
  const task = {id: 2, estimate_minutes: null} as Task;
  expect(taskEstimate(task, [])).toEqual({minutes: null, partial: false});
 });
});


describe('focus representation', () => {
 it('keeps the real running target when another task or completed step is requested', () => {
  const view = {server_now: '', selected_id: 99, sessions: [{...session, id: 99, task_id: 7, step_id: 1, state: 'paused' as const}, {...session, task_id: 2, step_id: 9}], steps: []};
  expect(focusTarget(view, {taskID: 7, stepID: 1})).toEqual({taskID: 2, stepID: 9, ready: false});
 });
 it('presents the next sibling stopped after current completes, and never completes the parent', () => {
  const steps = [{id: 9, task_id: 2, description: 'First', done: true, position: 0, estimate_minutes: null, deleted_at: null}, {id: 10, task_id: 2, description: 'Next', done: false, position: 1, estimate_minutes: null, deleted_at: null}];
  const view = {server_now: '', selected_id: 1, sessions: [{...session, task_id: 2, step_id: 9, state: 'ended' as const}], steps};
  expect(focusTarget(view, {taskID: 2, stepID: 9})).toEqual({taskID: 2, stepID: 10, ready: true});
  expect(view.sessions[0].state).toBe('ended');
  expect(focusTarget({...view, steps: steps.map((s) => ({...s, done: true}))}, null)).toEqual({taskID: 2, stepID: null, ready: true});
 });
 it('does not change the current action when an unrelated row is completed', () => {
  const view = {server_now: '', selected_id: null, sessions: [], steps: [{id: 9, task_id: 2, description: 'Other', done: true, position: 0, estimate_minutes: null, deleted_at: null}]};
  expect(focusTarget(view, {taskID: 2, stepID: 10})).toEqual({taskID: 2, stepID: 10, ready: false});
 });
});


describe('reorder sibling identity', () => {
 it('moves the dragged identity relative to siblings without mutating authoritative input', () => {
  const ids = [3, 7, 12, 19];
  expect(reorderStepIDs(ids, 3, 12)).toEqual([7, 12, 3, 19]);
  expect(reorderStepIDs(ids, 19, 7)).toEqual([3, 19, 7, 12]);
  expect(ids).toEqual([3, 7, 12, 19]);
 });
 it('produces no reorder write for canceled same or foreign targets', () => {
  expect(reorderStepIDs([3, 7], 3, 3)).toBeNull();
  expect(reorderStepIDs([3, 7], 3, 99)).toBeNull();
  expect(reorderStepIDs([3, 7], 99, 3)).toBeNull();
 });
});

describe('task estimate floor', () => {
 it('keeps the sum of all live steps including completed ones, and marks missing estimates', () => {
  const task = {id: 1, estimate_minutes: 5} as Task;
  const steps = [1,2,3].map((id) => ({id, task_id: 1, description: 'Action', estimate_minutes: 15, done: id === 1, position: id, deleted_at: null}));
  expect(taskEstimate(task, steps)).toEqual({minutes:45,partial:false});
  expect(taskEstimate({...task,estimate_minutes:60}, steps)).toEqual({minutes:60,partial:false});
  expect(taskEstimate(task, [...steps, {...steps[0],id:4,estimate_minutes:null}])).toEqual({minutes:45,partial:true});
 });
});
