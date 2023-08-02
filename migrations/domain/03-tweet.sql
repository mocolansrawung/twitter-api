DROP TABLE IF EXISTS `tweet`;

CREATE TABLE tweets (
    id CHAR(36) NOT NULL,
    content TEXT NOT NULL,
    retweets INT NOT NULL CHECK (retweets >= 0),
    created_at DATETIME NOT NULL,
    created_by CHAR(36) NOT NULL,
    updated_at DATETIME,
    updated_by CHAR(36),
    deleted_at DATETIME,
    deleted_by CHAR(36),
    PRIMARY KEY (id)
) ENGINE = InnoDB
DEFAULT CHARSET=utf8;


