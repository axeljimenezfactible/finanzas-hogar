# Extractos del contrato C1 · índice

Generado por `backend/internal/api/c1gen/build.py` (sin opciones o con `--extractos`); `build.py --check` los compara byte a byte con el generador. No se editan a mano.

**Regla (protocolo-agentes §4, K32):** estos archivos son los **únicos** archivos del contrato que leen los agentes de WF-1 y WF-3. `docs/contratos/openapi.yaml` es la fuente única y no se lee entero; si falta algo que no está en un extracto, se pide al integrador.

- Contrato C1 v1.5.1: 221 operaciones en 176 rutas y 267 esquemas.
- Cada operación aparece en exactamente un extracto. Cada extracto trae sus rutas y solo los esquemas que esas rutas usan, resueltos de forma transitiva y copiados en su `components`.
- Parámetros, encabezados, respuestas, seguridad y los esquemas `Error`, `Page`, `MoneyCents`, `Date`, `DateTime` (y los que ellos referencian) están en `componentes-comunes.yaml`; los extractos los citan con `$ref: componentes-comunes.yaml#/components/...`.
- Tope por archivo: 80 KB (80 000 bytes). Etiquetas divididas en partes consecutivas por orden de rutas: «Presupuesto» en `presupuesto-a.yaml`, `presupuesto-b.yaml`; «Compras y tiendas» en `compras-tiendas-a.yaml`, `compras-tiendas-b.yaml`; «Créditos» en `creditos-a.yaml`, `creditos-b.yaml`; «Ajustes» en `ajustes-a.yaml`, `ajustes-b.yaml`, `ajustes-c.yaml`, `ajustes-d.yaml`, `ajustes-e.yaml`.

| Archivo | Etiqueta | Operaciones | Rutas | Esquemas propios | Tamaño |
|---|---|---|---|---|---|
| `inicio.yaml` | Inicio | 15 | 15 | 20 | 66.4 KB (68038 bytes) |
| `gastos.yaml` | Gastos | 19 | 15 | 23 | 75.1 KB (76857 bytes) |
| `presupuesto-a.yaml` | Presupuesto | 13 | 11 | 14 | 45.5 KB (46580 bytes) |
| `presupuesto-b.yaml` | Presupuesto | 15 | 11 | 17 | 44.7 KB (45782 bytes) |
| `comidas.yaml` | Comidas | 16 | 12 | 19 | 48.2 KB (49399 bytes) |
| `compras-tiendas-a.yaml` | Compras y tiendas | 15 | 13 | 24 | 69.0 KB (70665 bytes) |
| `compras-tiendas-b.yaml` | Compras y tiendas | 16 | 13 | 25 | 62.6 KB (64081 bytes) |
| `creditos-a.yaml` | Créditos | 8 | 8 | 10 | 36.0 KB (36826 bytes) |
| `creditos-b.yaml` | Créditos | 11 | 8 | 23 | 62.9 KB (64376 bytes) |
| `salud.yaml` | Salud | 6 | 4 | 7 | 20.4 KB (20892 bytes) |
| `reportes.yaml` | Reportes | 10 | 10 | 21 | 58.4 KB (59768 bytes) |
| `ajustes-a.yaml` | Ajustes | 14 | 12 | 16 | 42.4 KB (43411 bytes) |
| `ajustes-b.yaml` | Ajustes | 14 | 13 | 21 | 52.6 KB (53910 bytes) |
| `ajustes-c.yaml` | Ajustes | 15 | 12 | 22 | 76.5 KB (78329 bytes) |
| `ajustes-d.yaml` | Ajustes | 14 | 13 | 21 | 56.7 KB (58063 bytes) |
| `ajustes-e.yaml` | Ajustes | 20 | 13 | 24 | 72.6 KB (74379 bytes) |
| `componentes-comunes.yaml` | (comunes) | — | — | 7 esquemas, 16 parámetros, 3 encabezados, 11 respuestas | 15.7 KB (16104 bytes) |
| **Total** | | **221** | | | |

## Operaciones por extracto

- `inicio.yaml`: `getHealth`, `getCurrentSession`, `search`, `listCalendarEvents`, `queryView`, `getDashboardMonthVsPlan`, `getDashboardFortnightFlow`, `getDashboardNextPayment`, `getDashboardWeekMeals`, `getDashboardBestBasket`, `getDashboardDebtCapacity`, `getDashboardMonthlyOverview`, `getBudgetSignals`, `getDashboardSavings`, `getPersonParticipation`
- `gastos.yaml`: `listExpenses`, `createExpense`, `getExpenseTotals`, `getExpenseSummary`, `getExpenseMonthlyTrend`, `getExpense`, `updateExpense`, `deleteExpense`, `listExpenseItems`, `duplicateExpense`, `markExpensePaid`, `uploadExpenseAttachment`, `downloadAttachment`, `deleteAttachment`, `createTicketScan`, `listTicketScans`, `getTicketScan`, `saveTicketScanDraft`, `confirmTicketScan`
- `presupuesto-a.yaml`: `listPeriods`, `listBudgets`, `createBudget`, `updateBudget`, `deleteBudget`, `getBudgetMonth`, `getBudgetByCategory`, `getBudgetFortnights`, `getAccruedGroceries`, `getBudgetAlerts`, `getBudgetMonthlyTrend`, `getCashFlowFortnights`, `getAccrualVsCash`
- `presupuesto-b.yaml`: `getCashFlowAssumptions`, `getLiveDebt`, `listIncomes`, `createIncome`, `getIncomeMonth`, `getIncomeMonthlyTrend`, `updateIncome`, `deleteIncome`, `getSavingsSummary`, `getSavingsMonthlyTrend`, `listSavingsGoals`, `listSavingsMovements`, `createSavingsMovement`, `updateSavingsMovement`, `deleteSavingsMovement`
- `comidas.yaml`: `listMealSlots`, `getMealPlanWeek`, `getMealPlanMonth`, `getMealPlanTrends`, `createMealPlanEntries`, `updateMealPlanEntry`, `deleteMealPlanEntry`, `copyPreviousWeek`, `listRecipes`, `createRecipe`, `suggestRecipeWithAi`, `getRecipe`, `updateRecipe`, `deleteRecipe`, `getRecipeUsage`, `duplicateRecipe`
- `compras-tiendas-a.yaml`: `getShoppingList`, `generateShoppingList`, `createShoppingListItem`, `updateShoppingListItem`, `assignShoppingListStore`, `listBaskets`, `getBasketComparison`, `listStores`, `listProductGroups`, `listProducts`, `createProduct`, `getProduct`, `updateProduct`, `createPresentation`, `updatePresentation`
- `compras-tiendas-b.yaml`: `listProductPrices`, `listNutritionFacts`, `createNutritionFact`, `getPriceSummary`, `createPrice`, `updatePrice`, `deletePrice`, `getExternalAgentReference`, `importExternalCatalog`, `importExternalPrices`, `importExternalRecipes`, `listInventoryItems`, `createInventoryItem`, `getInventorySummary`, `listInventoryMovements`, `adjustInventoryItem`
- `creditos-a.yaml`: `listCredits`, `getCreditTotals`, `getDebtTrend`, `getCredit`, `previewAvailability`, `getAvailabilitySeries`, `getCardStatus`, `listCreditTermRates`
- `creditos-b.yaml`: `listCreditPurchases`, `createCreditPurchase`, `previewCreditPurchase`, `getCreditPurchase`, `updateCreditPurchase`, `deleteCreditPurchase`, `listScheduledPayments`, `markScheduledPaymentPaid`, `deleteCreditPayment`, `listPromotions`, `compareScenarios`
- `salud.yaml`: `listHealthProfiles`, `getTodayKcal`, `listWeightEntries`, `createWeightEntry`, `updateWeightEntry`, `deleteWeightEntry`
- `reportes.yaml`: `listReports`, `generateMonthlyReport`, `generateExternalStatement`, `downloadReportFile`, `sendReport`, `listMonthlyReportMonths`, `getMonthlyReport`, `getExternalStatement`, `listNotificationDeliveries`, `listJobs`
- `ajustes-a.yaml`: `getOpenApiDocument`, `login`, `logout`, `changeOwnPassword`, `listUsers`, `listSessions`, `revokeSession`, `revokeSessions`, `getProjectionSettings`, `updateProjectionSettings`, `listPlannedItems`, `createPlannedItem`, `updatePlannedItem`, `previewPlannedItemsImport`
- `ajustes-b.yaml`: `importPlannedItems`, `createSavingsGoal`, `updateSavingsGoal`, `listCategories`, `createCategory`, `updateCategory`, `reorderCategories`, `listPersons`, `updatePerson`, `listCostShareRules`, `updateCostShareRule`, `createStore`, `updateStore`, `setStoreProductGroups`
- `ajustes-c.yaml`: `getPriceSources`, `getPriceSettings`, `updatePriceSettings`, `importProfeco`, `getRecipeSettings`, `reloadNutritionBase`, `updateInventoryItem`, `createCredit`, `updateCredit`, `createCreditTermRate`, `updateCreditTermRate`, `deleteCreditTermRate`, `createPromotion`, `updatePromotion`, `deletePromotion`
- `ajustes-d.yaml`: `getOptimizerSettings`, `updateOptimizerSettings`, `createHealthProfile`, `updateHealthProfile`, `updateJob`, `listNotifications`, `updateNotification`, `setNotificationRecipients`, `sendTestNotification`, `listNotificationRecipients`, `listTelegramChatCandidates`, `confirmTelegramChatCandidate`, `dismissTelegramChatCandidate`, `updateNotificationRecipient`
- `ajustes-e.yaml`: `listBackups`, `createBackup`, `restoreBackup`, `getSeedStatus`, `getSecuritySettings`, `updateSecuritySettings`, `getAlertSettings`, `updateAlertSettings`, `getDashboardSettings`, `updateDashboardSettings`, `getReportSettings`, `updateReportSettings`, `getDataSettings`, `updateDataSettings`, `getHealthSettings`, `updateHealthSettings`, `listAiProviders`, `updateAiProvider`, `getSystemStatus`, `getCaCertificate`
