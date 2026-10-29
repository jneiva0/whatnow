import { expect, setUp, startDay, test } from './fixtures';
import type { ExecutionState, Task } from '../src/lib/api';

test.beforeEach(async ({ request }) => {
 await setUp(request);
 await startDay(request, 'write the report and call mom');
});

test('steps can be created, edited, reordered, completed, reopened, split and undone', async ({page, request}) => {
 await page.goto('/');
 const card = page.locator('article', {hasText: 'write the report'});
 await card.getByLabel('More actions').click();
 await card.getByRole('button', {name: 'Edit', exact: true}).click();
 const sheet = page.getByRole('dialog');
 await sheet.getByLabel('Novo passo', {exact: true}).fill('Open the draft');
 await sheet.getByRole('button', {name: 'Adicionar passo'}).click();
 await expect(sheet.getByText('Open the draft', {exact: true})).toBeVisible();
 await sheet.getByLabel('Novo passo', {exact: true}).fill('Write one sentence');
 await sheet.getByRole('button', {name: 'Adicionar passo'}).click();
 await expect(sheet.getByText('Write one sentence', {exact: true})).toBeVisible();
 await sheet.getByRole('button', {name: 'Reordenar Write one sentence. Use setas para cima ou baixo.', exact: true}).focus();
 await page.keyboard.press('ArrowUp');
 await expect(sheet.locator('[data-step-row]').first()).toContainText('Write one sentence');
 await sheet.getByLabel('Concluir Write one sentence', {exact: true}).click();
 await expect(sheet.getByLabel('Reabrir Write one sentence')).toBeVisible();
 await sheet.getByLabel('Reabrir Write one sentence').click();
 const second = sheet.locator('[data-step-row]').filter({hasText: 'Open the draft'});
 await second.getByRole('button', {name: 'Open the draft', exact: true}).click();
 await sheet.getByLabel('Descrição do passo').fill('Open the document');
 await sheet.locator('.text-edit').getByRole('button', {name: 'Salvar', exact: true}).click();
 await sheet.getByRole('button', {name: 'Editar estimativa de Open the document', exact: true}).click();
 await sheet.getByLabel('Estimativa do passo em minutos').fill('2');
 await sheet.locator('.estimate-edit').getByRole('button', {name: 'Salvar', exact: true}).click();
 await expect(sheet.getByText('Open the document', {exact: true})).toBeVisible();
 await sheet.getByLabel('Mais ações para Write one sentence', {exact: true}).click();
 await sheet.getByRole('button', {name: 'Dividir este passo', exact: true}).click();
 await sheet.getByLabel('Substituir este passo por ações menores').fill('Type the title\nType the first sentence');
 await sheet.getByRole('button', {name: 'Dividir passo', exact: true}).click();
 await expect(sheet.getByText('Type the first sentence', {exact: true})).toBeVisible();
 await sheet.getByLabel('Mais ações para Open the document', {exact: true}).click();
 await sheet.getByRole('button', {name: 'Remover passo', exact: true}).click();
 await expect(sheet.getByText('Open the document', {exact: true})).toHaveCount(0);
 const state = await (await request.get('/api/v1/execution')).json() as ExecutionState;
 expect(state.steps.filter((s) => !s.deleted_at)).toHaveLength(2);
 const toast = page.getByRole('status').filter({hasText: 'Passo removido'});
 await toast.getByRole('button', {name: 'Undo', exact: true}).click();
 await expect(sheet.getByText('Open the document', {exact: true})).toBeVisible();
});

test('desktop and mobile share persistent focus, pause via SSE, correct time and complete explicitly', async ({page, request, browser}) => {
 const tasks = await (await request.get('/api/v1/tasks')).json() as Task[];
 const task = tasks.find((t) => t.title === 'write the report')!;
 const response = await request.post(`/api/v1/tasks/${task.id}/steps`, {data: {description: 'Open the draft', estimate_minutes: 3}});
 expect(response.ok()).toBe(true);
 await request.post(`/api/v1/tasks/${task.id}/steps`, {data: {description: 'Write one sentence', estimate_minutes: 2}});
 await page.goto('/');
 const dayCard = page.locator('article', {hasText: 'write the report'});
 await dayCard.locator('.card-steps summary').click();
 await dayCard.getByLabel('Concluir Open the draft', {exact: true}).click();
 await expect(dayCard.getByLabel('Reabrir Open the draft', {exact: true})).toBeVisible();
 expect((await (await request.get('/api/v1/execution')).json() as ExecutionState).sessions).toHaveLength(0);
 await dayCard.getByLabel('Reabrir Open the draft', {exact: true}).click();
 await dayCard.getByRole('button', {name: 'Começar Open the draft', exact: true}).click();

 await expect(page.getByRole('button', {name: 'Pausar foco', exact: true})).toBeVisible();
 const mobileContext = await browser.newContext({viewport: {width: 390, height: 844}, isMobile: true, hasTouch: true});
 const mobile = await mobileContext.newPage();
 try {
  await mobile.goto(`${new URL(page.url()).origin}/focus`);
  await expect(mobile.getByRole('heading', {name: 'Open the draft', exact: true})).toBeVisible();
  await mobile.getByRole('button', {name: 'Pausar foco', exact: true}).tap();
  await expect(page.getByText('Pausado', {exact: true})).toBeVisible();
  await page.reload();
  await expect(page.getByRole('button', {name: 'Iniciar ou retomar foco', exact: true})).toBeVisible();
  await page.getByText('Corrigir tempo registrado', {exact: true}).click();
  await page.getByRole('button', {name: 'Corrigir tempo', exact: true}).first().click();
  await page.getByLabel('Total correto em minutos').fill('0.5');
  await page.getByLabel('Motivo da correção').fill('Timer included time away');
  await page.getByRole('button', {name: 'Salvar correção', exact: true}).click();
  await expect(page.locator('.action-time')).toContainText('00:30');
  await expect(mobile.locator('.action-time')).toContainText('00:30');
  await mobile.getByRole('button', {name: 'Iniciar ou retomar foco', exact: true}).tap();
  await expect(page.getByRole('button', {name: 'Pausar foco', exact: true})).toBeVisible();
  await request.post('/api/v1/today/replan', {data: {}});
  await expect(page.getByRole('button', {name: 'Pausar foco', exact: true})).toBeVisible();
  await page.getByRole('button', {name: 'Concluir Open the draft', exact: true}).click();
  await expect(page.getByRole('heading', {name: 'Write one sentence', exact: true})).toBeVisible();
  await expect(page.getByText('Próximo passo · parado', {exact: true})).toBeVisible();
  expect((await (await request.get('/api/v1/execution')).json() as ExecutionState).sessions.some((s) => s.state === 'running')).toBe(false);
  expect((await (await request.get('/api/v1/tasks')).json() as Task[]).find((t) => t.id === task.id)?.status).toBe('open');
  await page.getByLabel('Concluir tarefa', {exact: true}).check();
  await expect(mobile.getByLabel('Reabrir tarefa', {exact: true})).toBeChecked();
  expect(await mobile.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true);
 } finally { await mobileContext.close(); }
});

test('API failure preserves the manual target and offers retry without pretending the timer started', async ({page, request}) => {
 const tasks = await (await request.get('/api/v1/tasks')).json() as Task[];
 await page.goto(`/focus?task=${tasks[0].id}`);
 await page.route('**/api/v1/focus/start', (route) => route.fulfill({status: 503, contentType: 'application/json', body: JSON.stringify({error: 'temporarily unavailable'})}));
 await page.getByRole('button', {name: 'Iniciar ou retomar foco', exact: true}).click();
 await expect(page.getByRole('alert')).toContainText('A ação não foi confirmada');
 await expect(page.getByRole('button', {name: 'Iniciar ou retomar foco', exact: true})).toBeVisible();
 await expect(page.getByRole('button', {name: 'Pausar foco', exact: true})).toHaveCount(0);
 await page.unroute('**/api/v1/focus/start');
 await page.getByRole('button', {name: 'Iniciar ou retomar foco', exact: true}).click();
 await expect(page.getByRole('button', {name: 'Pausar foco', exact: true})).toBeVisible();
 await page.getByRole('button', {name: 'Pausar foco', exact: true}).click();
 expect((await (await request.get('/api/v1/tasks')).json() as Task[])[0].status).toBe('open');
});
