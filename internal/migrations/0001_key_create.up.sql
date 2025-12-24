-- Таблица ссылок
CREATE TABLE links (
    lid SERIAL PRIMARY KEY,
    shortkey VARCHAR(16) NOT NULL UNIQUE,
    redirect TEXT NOT NULL,
    rentend TIMESTAMPTZ,
    isactive BOOLEAN DEFAULT TRUE
);
-- Таблица рефералов
CREATE TABLE referrals (
    rid SERIAL PRIMARY KEY,
    lid INTEGER NOT NULL,
    created TIMESTAMPTZ NOT NULL DEFAULT now(),
    useragent TEXT NOT NULL,
    CONSTRAINT fk_referrals_links FOREIGN KEY (lid) REFERENCES links (lid) ON UPDATE CASCADE ON DELETE CASCADE
);
-- Индексы для фильтрации аналитики
CREATE INDEX idx_referrals_lid ON referrals (lid);

CREATE INDEX idx_referrals_created ON referrals (created);