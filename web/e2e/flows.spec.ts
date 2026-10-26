import type { Page } from '@playwright/test';
import { expect, fake, setUp, startDay, test } from './fixtures';

const tier = (page: Page, name: string) =>
	page.locator('section').filter({ has: page.getByRole('heading', { name, exact: true }) });

test('a capture splits into tasks, and undo puts the note back', async ({ page, request }) => {
	await setUp(request);
	await startDay(request, '');
	await page.goto('/');

	await page.getByLabel('Capture').fill('buy milk and call mom');
	await page.getByLabel('Capture').press('Enter');

	const toast = page.getByRole('status').filter({ hasText: 'Split your note into 2 tasks' });
	await expect(toast).toBeVisible();
	await expect(page.getByText('buy milk', { exact: true })).toBeVisible();
	await expect(page.getByText('call mom', { exact: true })).toBeVisible();

	await toast.getByRole('button', { name: 'Undo' }).click();
	await expect(page.getByText('buy milk and call mom', { exact: true })).toBeVisible();
	await expect(page.getByText('call mom', { exact: true })).toHaveCount(0);
});

test('good morning: a dump becomes a plan in tiers, and Replan takes in a new task', async ({ page, request }) => {
	await setUp(request);
	await page.goto('/');
	await expect(page).toHaveURL(/\/morning$/);

	await page.getByLabel("What's on your plate").fill('buy milk, call mom, water the plants');
	await page.getByRole('button', { name: 'Next →' }).click();
	await page.getByRole('button', { name: 'pegando leve', exact: true }).click();
	await page.getByRole('button', { name: 'Plan my day →' }).click();

	await expect(page.getByText(fake.opening)).toBeVisible();
	await page.getByRole('button', { name: "Let's go →" }).click();
	await expect(page).toHaveURL(/\/$/);
	await expect(tier(page, 'Must (Essencial)').getByText('buy milk', { exact: true })).toBeVisible();
	await expect(tier(page, 'Should (Importante)').getByText('call mom', { exact: true })).toBeVisible();
	await expect(tier(page, 'Could (Se der)').getByText('water the plants', { exact: true })).toBeVisible();

	await page.getByLabel('Capture').fill('feed the cat');
	await page.getByLabel('Capture').press('Enter');
	await expect(tier(page, 'Unplanned').getByText('feed the cat', { exact: true })).toBeVisible();
	await page.getByRole('button', { name: 'Replan' }).click();
	await expect(tier(page, 'Could (Se der)').getByText('feed the cat', { exact: true })).toBeVisible();
});

test('what now picks the first free task, skips a pass, and opens execution', async ({ page, request }) => {
	await setUp(request);
	await startDay(request, 'buy milk and call mom');
	await page.goto('/');

	await page.getByRole('button', { name: 'E agora?' }).click();
	const pick = page.locator('section.whatnow');
	await expect(pick.getByRole('heading', { name: 'buy milk' })).toBeVisible();
	await expect(pick.getByText(fake.reason)).toBeVisible();

	await pick.getByRole('button', { name: 'Not this one' }).click();
	await pick.getByRole('button', { name: 'Ask again' }).click();
	await expect(pick.getByRole('heading', { name: 'call mom' })).toBeVisible();

	await pick.getByRole('button', { name: 'Começar' }).click();
	await expect(page.getByRole('heading', { name: 'call mom', exact: true })).toBeVisible();
	await page.getByRole('button', { name: 'Iniciar foco', exact: true }).click();
	await page.getByRole('button', { name: 'Concluir tarefa', exact: true }).click();
	await expect(page.getByText('Ação concluída', { exact: true })).toBeVisible();
});
