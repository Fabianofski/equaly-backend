CREATE TABLE ExpenseLists
(
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    color        VARCHAR(10),
    emoji        VARCHAR(10),
    title        VARCHAR(50),
    creatorId    VARCHAR(50),
    currency     VARCHAR(3),
    participants json,
    inviteCode   VARCHAR(10) DEFAULT '1234567890'
);

CREATE TABLE Expenses
(
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    expenseListId UUID,
    buyer         VARCHAR(50),
    amount        DECIMAL(10, 2),
    description   VARCHAR(255),
    participants  TEXT,
    date          TIMESTAMPTZ,
    FOREIGN KEY (expenseListId) REFERENCES ExpenseLists (id)
);

CREATE TABLE ExpenseListsUsers(
    expenseListId UUID,
    userId VARCHAR(50),
    FOREIGN KEY (expenseListId) REFERENCES ExpenseLists (id)
);
