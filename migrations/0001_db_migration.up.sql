CREATE DATABASE News;

CREATE TABLE News (
  Id BIGINT PRIMARY KEY NOT NULL,
  Title TEXT NOT NULL,
  Content TEXT NOT NULL
);

CREATE TABLE NewsCategories (
  NewsId BIGINT NOT NULL,
  CategoryId BIGINT NOT NULL,
  PRIMARY KEY (NewsId, CategoryId),
  FOREIGN KEY (NewsId) REFERENCES News(Id) ON DELETE CASCADE
);

INSERT INTO
    News
VALUES
    (1, 'FirstTitle', 'Content News'),
    (64, 'Lorem ipsum', 'Content!'),
    (23, 'Third', '9 is a lucky number');

INSERT INTO 
    NewsCategories
VALUES
    (1, 1),
    (1, 2),
    (1, 5),
    (64, 1),
    (64, 3),
    (64, 9),
    (23, 5),
    (23, 2),
    (23, 3);
