# Extractos del contrato C1 · índice

Generado por `backend/internal/api/c1gen/build.py` (sin opciones o con `--extractos`); `build.py --check` los compara byte a byte con el generador. No se editan a mano.

**Regla (protocolo-agentes §4, K32):** estos archivos son los **únicos** archivos del contrato que leen los agentes de WF-1 y WF-3. `docs/contratos/openapi.yaml` es la fuente única y no se lee entero; si falta algo que no está en un extracto, se pide al integrador.

- Contrato C1 v1.4.0: 216 operaciones en 172 rutas y 262 esquemas.
- Cada operación aparece en exactamente un extracto. Cada extracto trae sus rutas y solo los esquemas que esas rutas usan, resueltos de forma transitiva y copiados en su `components`.
- Parámetros, encabezados, respuestas, seguridad y los esquemas `Error`, `Page`, `MoneyCents`, `Date`, `DateTime` (y los que ellos referencian) están en `componentes-comunes.yaml`; los extractos los citan con `$ref: componentes-comunes.yaml#/components/...`.
- Tope por archivo: 80 KB (80 000 bytes). Etiquetas divididas en partes consecutivas por orden de rutas: «Presupuesto» en `presupuesto-a.yaml`, `presupuesto-b.yaml`; «Compras y tiendas» en `compras-tiendas-a.yaml`, `compras-tiendas-b.yaml`; «Créditos» en `creditos-a.yaml`, `creditos-b.yaml`; «Ajustes» en `ajustes-a.yaml`, `ajustes-b.yaml`, `ajustes-c.yaml`, `ajustes-d.yaml`, `ajustes-e.yaml`, `ajustes-f.yaml`.

| Archivo | Etiqueta | Operaciones | Rutas | Esquemas propios | Tamaño |
|---|---|---|---|---|---|
| `inicio.yaml` | Inicio | 12 | 12 | 17 | 47.7 KB (48845 bytes) |
| `gastos.yaml` | Gastos | 19 | 15 | 23 | 75.1 KB (76857 bytes) |
| `presupuesto-a.yaml` | Presupuesto | 13 | 11 | 14 | 45.5 KB (46580 bytes) |
| `presupuesto-b.yaml` | Presupuesto | 15 | 11 | 17 | 44.0 KB (45012 bytes) |
| `comidas.yaml` | Comidas | 16 | 12 | 19 | 48.2 KB (49399 bytes) |
| `compras-tiendas-a.yaml` | Compras y tiendas | 15 | 13 | 24 | 69.0 KB (70665 bytes) |
| `compras-tiendas-b.yaml` | Compras y tiendas | 16 | 13 | 25 | 62.6 KB (64081 bytes) |
| `creditos-a.yaml` | Créditos | 8 | 8 | 10 | 34.4 KB (35245 bytes) |
| `creditos-b.yaml` | Créditos | 11 | 8 | 23 | 62.9 KB (64376 bytes) |
| `salud.yaml` | Salud | 6 | 4 | 7 | 20.4 KB (20892 bytes) |
| `reportes.yaml` | Reportes | 10 | 10 | 21 | 57.2 KB (58550 bytes) |
| `ajustes-a.yaml` | Ajustes | 12 | 10 | 14 | 34.1 KB (34933 bytes) |
| `ajustes-b.yaml` | Ajustes | 11 | 10 | 16 | 45.8 KB (46871 bytes) |
| `ajustes-c.yaml` | Ajustes | 12 | 11 | 22 | 53.0 KB (54267 bytes) |
| `ajustes-d.yaml` | Ajustes | 13 | 10 | 19 | 77.3 KB (79160 bytes) |
| `ajustes-e.yaml` | Ajustes | 11 | 10 | 14 | 32.1 KB (32876 bytes) |
| `ajustes-f.yaml` | Ajustes | 16 | 11 | 20 | 56.7 KB (58046 bytes) |
| `componentes-comunes.yaml` | (comunes) | — | — | 7 esquemas, 14 parámetros, 3 encabezados, 11 respuestas | 15.0 KB (15340 bytes) |
| **Total** | | **216** | | | |

## Operaciones por extracto

- `inicio.yaml`: `getHealth`, `getCurrentSession`, `search`, `listCalendarEvents`, `queryView`, `getDashboardMonthVsPlan`, `getDashboardFortnightFlow`, `getDashboardNextPayment`, `getDashboardWeekMeals`, `getDashboardBestBasket`, `getDashboardSavings`, `getPersonParticipation`
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
- `ajustes-a.yaml`: `getOpenApiDocument`, `login`, `logout`, `changeOwnPassword`, `listUsers`, `listSessions`, `revokeSession`, `revokeSessions`, `getProjectionSettings`, `updateProjectionSettings`, `listPlannedItems`, `createPlannedItem`
- `ajustes-b.yaml`: `updatePlannedItem`, `previewPlannedItemsImport`, `importPlannedItems`, `createSavingsGoal`, `updateSavingsGoal`, `listCategories`, `createCategory`, `updateCategory`, `reorderCategories`, `listPersons`, `updatePerson`
- `ajustes-c.yaml`: `listCostShareRules`, `updateCostShareRule`, `createStore`, `updateStore`, `setStoreProductGroups`, `getPriceSources`, `getPriceSettings`, `updatePriceSettings`, `importProfeco`, `getRecipeSettings`, `reloadNutritionBase`, `updateInventoryItem`
- `ajustes-d.yaml`: `createCredit`, `updateCredit`, `createCreditTermRate`, `updateCreditTermRate`, `deleteCreditTermRate`, `createPromotion`, `updatePromotion`, `deletePromotion`, `getOptimizerSettings`, `updateOptimizerSettings`, `createHealthProfile`, `updateHealthProfile`, `updateJob`
- `ajustes-e.yaml`: `listNotifications`, `updateNotification`, `setNotificationRecipients`, `sendTestNotification`, `listNotificationRecipients`, `listTelegramChatCandidates`, `confirmTelegramChatCandidate`, `dismissTelegramChatCandidate`, `updateNotificationRecipient`, `listBackups`, `createBackup`
- `ajustes-f.yaml`: `restoreBackup`, `getSeedStatus`, `getSecuritySettings`, `updateSecuritySettings`, `getAlertSettings`, `updateAlertSettings`, `getReportSettings`, `updateReportSettings`, `getDataSettings`, `updateDataSettings`, `getHealthSettings`, `updateHealthSettings`, `listAiProviders`, `updateAiProvider`, `getSystemStatus`, `getCaCertificate`
