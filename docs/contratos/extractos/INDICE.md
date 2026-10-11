# Extractos del contrato C1 · índice

Generado por `backend/internal/api/c1gen/build.py` (sin opciones o con `--extractos`); `build.py --check` los compara byte a byte con el generador. No se editan a mano.

**Regla (protocolo-agentes §4, K32):** estos archivos son los **únicos** archivos del contrato que leen los agentes de WF-1 y WF-3. `docs/contratos/openapi.yaml` es la fuente única y no se lee entero; si falta algo que no está en un extracto, se pide al integrador.

- Contrato C1 v1.1.0: 201 operaciones en 162 rutas y 242 esquemas.
- Cada operación aparece en exactamente un extracto. Cada extracto trae sus rutas y solo los esquemas que esas rutas usan, resueltos de forma transitiva y copiados en su `components`.
- Parámetros, encabezados, respuestas, seguridad y los esquemas `Error`, `Page`, `MoneyCents`, `Date`, `DateTime` (y los que ellos referencian) están en `componentes-comunes.yaml`; los extractos los citan con `$ref: componentes-comunes.yaml#/components/...`.
- Tope por archivo: 80 KB (80 000 bytes). Etiquetas divididas en partes consecutivas por orden de rutas: «Presupuesto» en `presupuesto-a.yaml`, `presupuesto-b.yaml`; «Compras y tiendas» en `compras-tiendas-a.yaml`, `compras-tiendas-b.yaml`; «Créditos» en `creditos-a.yaml`, `creditos-b.yaml`; «Ajustes» en `ajustes-a.yaml`, `ajustes-b.yaml`, `ajustes-c.yaml`, `ajustes-d.yaml`, `ajustes-e.yaml`.

| Archivo | Etiqueta | Operaciones | Rutas | Esquemas propios | Tamaño |
|---|---|---|---|---|---|
| `inicio.yaml` | Inicio | 12 | 12 | 17 | 46.3 KB (47373 bytes) |
| `gastos.yaml` | Gastos | 19 | 15 | 23 | 74.4 KB (76202 bytes) |
| `presupuesto-a.yaml` | Presupuesto | 13 | 11 | 14 | 42.9 KB (43888 bytes) |
| `presupuesto-b.yaml` | Presupuesto | 15 | 11 | 17 | 42.7 KB (43754 bytes) |
| `comidas.yaml` | Comidas | 16 | 12 | 19 | 47.4 KB (48574 bytes) |
| `compras-tiendas-a.yaml` | Compras y tiendas | 13 | 11 | 23 | 60.3 KB (61703 bytes) |
| `compras-tiendas-b.yaml` | Compras y tiendas | 14 | 11 | 19 | 47.5 KB (48594 bytes) |
| `creditos-a.yaml` | Créditos | 8 | 8 | 10 | 32.4 KB (33140 bytes) |
| `creditos-b.yaml` | Créditos | 11 | 8 | 23 | 61.8 KB (63266 bytes) |
| `salud.yaml` | Salud | 6 | 4 | 7 | 18.1 KB (18563 bytes) |
| `reportes.yaml` | Reportes | 10 | 10 | 21 | 56.1 KB (57440 bytes) |
| `ajustes-a.yaml` | Ajustes | 13 | 11 | 15 | 38.4 KB (39271 bytes) |
| `ajustes-b.yaml` | Ajustes | 12 | 11 | 20 | 43.6 KB (44597 bytes) |
| `ajustes-c.yaml` | Ajustes | 12 | 11 | 19 | 55.7 KB (57015 bytes) |
| `ajustes-d.yaml` | Ajustes | 14 | 11 | 18 | 50.9 KB (52118 bytes) |
| `ajustes-e.yaml` | Ajustes | 13 | 12 | 17 | 31.9 KB (32675 bytes) |
| `componentes-comunes.yaml` | (comunes) | — | — | 7 esquemas, 14 parámetros, 3 encabezados, 10 respuestas | 14.0 KB (14322 bytes) |
| **Total** | | **201** | | | |

## Operaciones por extracto

- `inicio.yaml`: `getHealth`, `getCurrentSession`, `search`, `listCalendarEvents`, `queryView`, `getDashboardMonthVsPlan`, `getDashboardFortnightFlow`, `getDashboardNextPayment`, `getDashboardWeekMeals`, `getDashboardBestBasket`, `getDashboardSavings`, `getPersonParticipation`
- `gastos.yaml`: `listExpenses`, `createExpense`, `getExpenseTotals`, `getExpenseSummary`, `getExpenseMonthlyTrend`, `getExpense`, `updateExpense`, `deleteExpense`, `listExpenseItems`, `duplicateExpense`, `markExpensePaid`, `uploadExpenseAttachment`, `downloadAttachment`, `deleteAttachment`, `createTicketScan`, `listTicketScans`, `getTicketScan`, `saveTicketScanDraft`, `confirmTicketScan`
- `presupuesto-a.yaml`: `listPeriods`, `listBudgets`, `createBudget`, `updateBudget`, `deleteBudget`, `getBudgetMonth`, `getBudgetByCategory`, `getBudgetFortnights`, `getAccruedGroceries`, `getBudgetAlerts`, `getBudgetMonthlyTrend`, `getCashFlowFortnights`, `getAccrualVsCash`
- `presupuesto-b.yaml`: `getCashFlowAssumptions`, `getLiveDebt`, `listIncomes`, `createIncome`, `getIncomeMonth`, `getIncomeMonthlyTrend`, `updateIncome`, `deleteIncome`, `getSavingsSummary`, `getSavingsMonthlyTrend`, `listSavingsGoals`, `listSavingsMovements`, `createSavingsMovement`, `updateSavingsMovement`, `deleteSavingsMovement`
- `comidas.yaml`: `listMealSlots`, `getMealPlanWeek`, `getMealPlanMonth`, `getMealPlanTrends`, `createMealPlanEntries`, `updateMealPlanEntry`, `deleteMealPlanEntry`, `copyPreviousWeek`, `listRecipes`, `createRecipe`, `suggestRecipeWithAi`, `getRecipe`, `updateRecipe`, `deleteRecipe`, `getRecipeUsage`, `duplicateRecipe`
- `compras-tiendas-a.yaml`: `getShoppingList`, `generateShoppingList`, `createShoppingListItem`, `updateShoppingListItem`, `assignShoppingListStore`, `listBaskets`, `getBasketComparison`, `listStores`, `listProductGroups`, `listProducts`, `createProduct`, `getProduct`, `updateProduct`
- `compras-tiendas-b.yaml`: `createPresentation`, `updatePresentation`, `listProductPrices`, `listNutritionFacts`, `createNutritionFact`, `getPriceSummary`, `createPrice`, `updatePrice`, `deletePrice`, `listInventoryItems`, `createInventoryItem`, `getInventorySummary`, `listInventoryMovements`, `adjustInventoryItem`
- `creditos-a.yaml`: `listCredits`, `getCreditTotals`, `getDebtTrend`, `getCredit`, `previewAvailability`, `getAvailabilitySeries`, `getCardStatus`, `listCreditTermRates`
- `creditos-b.yaml`: `listCreditPurchases`, `createCreditPurchase`, `previewCreditPurchase`, `getCreditPurchase`, `updateCreditPurchase`, `deleteCreditPurchase`, `listScheduledPayments`, `markScheduledPaymentPaid`, `deleteCreditPayment`, `listPromotions`, `compareScenarios`
- `salud.yaml`: `listHealthProfiles`, `getTodayKcal`, `listWeightEntries`, `createWeightEntry`, `updateWeightEntry`, `deleteWeightEntry`
- `reportes.yaml`: `listReports`, `generateMonthlyReport`, `generateExternalStatement`, `downloadReportFile`, `sendReport`, `listMonthlyReportMonths`, `getMonthlyReport`, `getExternalStatement`, `listNotificationDeliveries`, `listJobs`
- `ajustes-a.yaml`: `getOpenApiDocument`, `login`, `logout`, `changeOwnPassword`, `listUsers`, `listSessions`, `revokeSession`, `revokeSessions`, `getProjectionSettings`, `updateProjectionSettings`, `listPlannedItems`, `createPlannedItem`, `updatePlannedItem`
- `ajustes-b.yaml`: `previewPlannedItemsImport`, `importPlannedItems`, `createSavingsGoal`, `updateSavingsGoal`, `listCategories`, `createCategory`, `updateCategory`, `reorderCategories`, `listPersons`, `updatePerson`, `listCostShareRules`, `updateCostShareRule`
- `ajustes-c.yaml`: `createStore`, `updateStore`, `setStoreProductGroups`, `getPriceSources`, `getPriceSettings`, `updatePriceSettings`, `importProfeco`, `getRecipeSettings`, `reloadNutritionBase`, `updateInventoryItem`, `createCredit`, `updateCredit`
- `ajustes-d.yaml`: `createCreditTermRate`, `updateCreditTermRate`, `deleteCreditTermRate`, `createPromotion`, `updatePromotion`, `deletePromotion`, `getOptimizerSettings`, `updateOptimizerSettings`, `createHealthProfile`, `updateHealthProfile`, `updateJob`, `listNotifications`, `updateNotification`, `setNotificationRecipients`
- `ajustes-e.yaml`: `sendTestNotification`, `listNotificationRecipients`, `listTelegramChatCandidates`, `confirmTelegramChatCandidate`, `updateNotificationRecipient`, `listBackups`, `createBackup`, `restoreBackup`, `getSeedStatus`, `listAiProviders`, `updateAiProvider`, `getSystemStatus`, `getCaCertificate`
