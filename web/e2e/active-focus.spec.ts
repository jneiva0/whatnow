import { expect, setUp, startDay, test } from './fixtures';
import type { Task } from '../src/lib/api';

test('active focus stays visible on the day and pauses the existing session', async ({ page, request }) => {
	await setUp(request);
	await startDay(request, 'write the report and call mom');
	const tasks = await (await request.get('/api/v1/tasks')).json() as Task[];
	const task = tasks.find((task) => task.title === 'write the report')!;
	await request.post(`/api/v1/tasks/${task.id}/steps`, { data: { description: 'Open the draft', estimate_minutes: 3 } });
	await page.goto('/?view=actions');
	await page.getByRole('region', { name: 'Sugestão para começar' }).getByRole('button', { name: 'Começar', exact: true }).click();
	await expect(page).toHaveURL(/\/focus\?/);
	await expect(page.getByRole('complementary', { name: 'Foco em andamento' })).toHaveCount(0);
	await page.goto('/?view=actions');
	const bar = page.getByRole('complementary', { name: 'Foco em andamento' });
	await expect(bar).toContainText('Open the draft');
	await expect(bar.locator('.clock')).toHaveAttribute('aria-label', /^Tempo restante:/);
	await expect(bar.getByLabel('Tempo dedicado / estimativa', { exact: true })).toContainText('/ 03:00');
	await expect(bar.getByRole('progressbar')).toBeVisible();
	await page.getByRole('button', { name: 'Fechar chat', exact: true }).click();
	const card = page.locator(`article[data-task="${task.id}"]`);
	await card.getByRole('combobox', { name: /Energia da tarefa/ }).click();
	await page.getByRole('option', { name: 'Baixa', exact: true }).click();
	await expect(card.getByRole('combobox', { name: 'Energia da tarefa: Baixa' })).toBeVisible();
	await card.getByRole('button', { name: 'Editar estimativa da tarefa em minutos' }).click();
	await card.getByRole('textbox', { name: 'Editar estimativa da tarefa em minutos' }).fill('5');
	await page.keyboard.press('Enter');
	await expect(card.getByRole('button', { name: 'Editar estimativa da tarefa em minutos' })).toHaveText('5 min');
	for (const width of [1440, 390, 320]) {
		await page.setViewportSize({ width, height: 700 });
		await expect(bar).toBeVisible();
		expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true);
		await page.evaluate(() => window.scrollTo(0, document.body.scrollHeight));
		expect((await bar.boundingBox())!.y).toBe(0);
		await page.screenshot({ path: `/tmp/whatnow-focus-${width}.png` });
	}
	await bar.getByRole('link').click();
	await expect(page).toHaveURL(/\/focus\?/);
	await page.goto('/');
	await bar.getByRole('button', { name: 'Pausar Open the draft', exact: true }).click();
	await expect(bar).toHaveCount(0);
	const state = await (await request.get('/api/v1/execution')).json();
	expect(state.sessions).toHaveLength(1);
	expect(state.sessions[0].state).toBe('paused');
	await card.getByRole('button', { name: 'Editar estimativa da tarefa em minutos' }).click();
	await card.getByRole('textbox', { name: 'Editar estimativa da tarefa em minutos' }).fill('28');
	await page.keyboard.press('Enter');
	await expect(card.getByRole('button', { name: 'Editar estimativa da tarefa em minutos' })).toHaveText('28 min');
	for (const width of [1440, 390, 320]) {
		await page.setViewportSize({ width, height: 800 });
		await page.evaluate(() => window.scrollTo(0, 0));
		expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true);
		await expect(card.locator('.summary-total')).toContainText('28:00');
		const estimate = card.getByRole('button', { name: 'Editar estimativa da tarefa em minutos' });
		expect(await estimate.evaluate((node) => getComputedStyle(node).whiteSpace)).toBe('nowrap');
		await page.screenshot({ path: `/tmp/whatnow-compact-card-${width}.png`, fullPage: true });
	}
});
