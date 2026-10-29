import { expect, setUp, startDay, test } from './fixtures';
import type { Task } from '../src/lib/api';

test('Actions suggests unplanned work, refreshes after capture, clarifies in chat and respects focus', async ({page, request}) => {
 await setUp(request); await startDay(request, 'write report and call mom');
 const create = await request.post('/api/v1/tasks', {data: {title: 'Ritual para começar', notes: 'Beber água primeiro'}});
 const ritual = (await create.json()).data as Task;
 const stepResponse = await request.post(`/api/v1/tasks/${ritual.id}/steps`, {data: {description: 'Beber água', estimate_minutes: 1}});
 const step = (await stepResponse.json()).data;
 const nextResponse = await request.post(`/api/v1/tasks/${ritual.id}/steps`, {data: {description: 'Ir ao banheiro', estimate_minutes: 2}});
 const next = (await nextResponse.json()).data; let selectedStep = step;
 let requests = 0; let recentContext = '';
 await page.route('**/api/v1/whatnow', async (route) => {
  requests++; recentContext = route.request().postDataJSON().context;
  await route.fulfill({json: {task_id: ritual.id, step_id: selectedStep.id, reason: 'Primeiro movimento do ritual que você descreveu.', fallback: false}});
 });
 await page.goto('/?view=actions');
 const suggestion = page.getByRole('region', {name: 'Sugestão para começar'});
 await expect(suggestion.getByRole('heading', {name: 'Beber água', exact: true})).toBeVisible();
 await expect(page.locator(`article[data-task="${ritual.id}"]`)).toHaveCount(0);
 await page.getByRole('link', {name: 'Prioridades', exact: true}).click();
 await page.getByRole('link', {name: 'Ações', exact: true}).click();
 expect(requests).toBe(1);
 await suggestion.getByRole('button', {name: 'Ajustar contexto'}).click();
 await page.getByRole('textbox', {name: 'Capture', exact: true}).fill('take a walk'); await page.getByRole('textbox', {name: 'Capture', exact: true}).press('Enter');
 await expect(page.getByRole('status').filter({hasText: 'Tidying'})).toHaveCount(0);
 await expect(suggestion.getByRole('textbox', {name: 'O que considerar agora?'})).toBeVisible();
 expect(requests).toBe(1);
 await suggestion.getByRole('button', {name: 'Ajustar contexto'}).click();
 await expect.poll(() => requests).toBe(2);
 expect(recentContext).toContain('take a walk');
 await suggestion.getByRole('button', {name: 'Me ajude a organizar'}).click();
 await expect(page.getByText(/Me ajude a esclarecer a tarefa/)).toBeVisible();
 await expect(page.getByText('A fake reply. Nothing changed.', {exact: true})).toBeVisible();
 await page.getByRole('button', {name: 'Fechar chat'}).click();
 for (const width of [1440, 390, 320]) {
  await page.setViewportSize({width, height: 800});
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true);
  await page.screenshot({path: `/tmp/whatnow-suggestion-${width}.png`});
 }
 await suggestion.getByRole('button', {name: 'Começar', exact: true}).click();
 await expect(page).toHaveURL(/\/focus\?/);
 await page.goto('/?view=actions');
 await expect(page.getByRole('heading', {name: 'Continue de onde está'})).toBeVisible();
 await expect(suggestion).toHaveCount(0);
 const before = requests;
 await expect(page.getByRole('link', {name: 'Voltar ao foco', exact: true})).toBeVisible();
 expect(requests).toBe(before);
 selectedStep = next;
 const finish = await request.put(`/api/v1/tasks/${ritual.id}/steps/${step.id}`, {data: {description: 'Beber água', estimate_minutes: 1, done: true}});
 expect(finish.ok()).toBe(true);
 await expect(suggestion.getByRole('heading', {name: 'Ir ao banheiro', exact: true})).toBeVisible();
 const execution = await (await request.get('/api/v1/execution')).json();
 expect(execution.sessions.some((s: {state: string}) => s.state === 'running')).toBe(false);
});

test('temporary refusal uses context, retains tasks, and transport failure leaves manual choices', async ({page, request}) => {
 await setUp(request); await startDay(request, 'first and second');
 const tasks = await (await request.get('/api/v1/tasks')).json() as Task[];
 let calls = 0; let body: Record<string, unknown> = {};
 await page.route('**/api/v1/whatnow', async (route) => {
  body = route.request().postDataJSON(); calls++;
  if (calls === 3) { await route.abort(); return; }
  await route.fulfill({json: {task_id: tasks[calls - 1].id, step_id: null, reason: 'Concrete task', fallback: calls === 2}});
 });
 await page.goto('/?view=actions');
 const panel = page.getByRole('region', {name: 'Sugestão para começar'});
 await expect(panel.getByRole('heading', {name: tasks[0].title, exact: true})).toBeVisible();
 await panel.getByRole('button', {name: 'Outra sugestão'}).click();
 await expect(panel.getByRole('heading', {name: tasks[1].title, exact: true})).toBeVisible();
 expect(body.reason).toBe(''); expect(body.exclude_task_id).toBe(tasks[0].id);
 await expect(panel.getByText(/Sugestão básica por prazo/)).toBeVisible();
 await panel.getByRole('button', {name: 'Ajustar contexto'}).click();
 await panel.getByRole('textbox', {name: 'O que considerar agora?'}).fill('Estou cansado agora');
 await panel.getByRole('button', {name: 'Usar este contexto'}).click();
 await expect(panel.getByRole('alert')).toContainText('Não foi possível');
 expect(body.context).toBe('Estou cansado agora'); expect(body.reason).toBe('');
 expect(await (await request.get('/api/v1/tasks')).json()).toEqual(tasks);
 await expect(page.locator('article').first()).toBeVisible(); expect(calls).toBe(3);
});
