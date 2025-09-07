-- Companies テーブル作成
CREATE TABLE companies (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(20),
    address VARCHAR(500),
    website VARCHAR(255),
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- インデックス作成
CREATE INDEX idx_companies_email ON companies(email);
CREATE INDEX idx_companies_name ON companies(name);
CREATE INDEX idx_companies_created_at ON companies(created_at);

-- updated_at 自動更新トリガー
CREATE TRIGGER update_companies_updated_at
    BEFORE UPDATE ON companies
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Company-User 中間テーブル作成（多対多関係）
CREATE TABLE company_users (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    company_id INTEGER NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    role VARCHAR(50) DEFAULT 'member',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, company_id)
);

-- company_users インデックス作成
CREATE INDEX idx_company_users_user_id ON company_users(user_id);
CREATE INDEX idx_company_users_company_id ON company_users(company_id);
CREATE INDEX idx_company_users_role ON company_users(role);

-- company_users updated_at 自動更新トリガー
CREATE TRIGGER update_company_users_updated_at
    BEFORE UPDATE ON company_users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- サンプルデータ挿入
INSERT INTO companies (name, email, phone, address, website, description) VALUES
    ('株式会社サンプル', 'info@sample.co.jp', '03-1234-5678', '東京都渋谷区サンプル町1-1-1', 'https://sample.co.jp', 'IT関連のサービスを提供する会社です。'),
    ('テスト商事株式会社', 'contact@test-corp.jp', '06-9876-5432', '大阪府大阪市テスト区2-2-2', 'https://test-corp.jp', '商社として様々な商品を扱っています。'),
    ('エンジニアリング合同会社', 'hello@engineering.com', '045-1111-2222', '神奈川県横浜市エンジニア区3-3-3', 'https://engineering.com', 'エンジニアリングソリューションを提供しています。');

-- サンプル関係データ挿入
INSERT INTO company_users (user_id, company_id, role) VALUES
    (1, 1, 'admin'),    -- 山田太郎 → 株式会社サンプル（管理者）
    (2, 1, 'member'),   -- 佐藤花子 → 株式会社サンプル（メンバー）
    (2, 2, 'admin'),    -- 佐藤花子 → テスト商事（管理者）
    (3, 3, 'admin');    -- 田中一郎 → エンジニアリング合同会社（管理者）

-- テーブルコメント
COMMENT ON TABLE companies IS '会社情報テーブル';
COMMENT ON COLUMN companies.id IS '会社ID（主キー）';
COMMENT ON COLUMN companies.name IS '会社名';
COMMENT ON COLUMN companies.email IS 'メールアドレス（ユニーク）';
COMMENT ON COLUMN companies.phone IS '電話番号';
COMMENT ON COLUMN companies.address IS '住所';
COMMENT ON COLUMN companies.website IS 'ウェブサイトURL';
COMMENT ON COLUMN companies.description IS '会社説明';
COMMENT ON COLUMN companies.created_at IS '作成日時';
COMMENT ON COLUMN companies.updated_at IS '更新日時';

COMMENT ON TABLE company_users IS 'ユーザー-会社関係テーブル';
COMMENT ON COLUMN company_users.id IS '関係ID（主キー）';
COMMENT ON COLUMN company_users.user_id IS 'ユーザーID（外部キー）';
COMMENT ON COLUMN company_users.company_id IS '会社ID（外部キー）';
COMMENT ON COLUMN company_users.role IS 'ユーザーの役割（admin, member）';
COMMENT ON COLUMN company_users.created_at IS '関係作成日時';
COMMENT ON COLUMN company_users.updated_at IS '関係更新日時';