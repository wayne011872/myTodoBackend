CREATE TABLE IF NOT EXISTS todoitem (
    id bigserial PRIMARY KEY,            -- 編號
    title varchar(255) NOT NULL,    -- 標題
    detail varchar(2000),      -- 內容
    completed boolean NOT NULL,      -- 完成狀態
    startTime timestamp NOT NULL,             -- 開始時間
    endTime timestamp NOT NULL,               -- 結束時間
    createdTime timestamp DEFAULT NOW(),        -- 建立時間
    updatedTime timestamp                       -- 更新時間
);