package constants

var SQLStatements = map[string]string{
	"books": `SELECT
				cards.id card_id,
				cards.name card_name,
				cards.lastname card_lastname,
				cards.surname card_surname,
				cards.phone card_phone,
				schools.id school_id,
				schools.name school_name,
				books.id book_id,
				books.name book_name,
				books.author book_author,
				books.book_id book_ownID
			FROM
				books
			LEFT JOIN cards_books
				ON books.id = cards_books.book_id
			INNER JOIN cards
				ON cards.id = cards_books.card_id
			LEFT JOIN schools
				ON schools.id = cards.school_id`,
	"cards": `SELECT
				cards.id card_id,
				cards.name card_name,
				cards.lastname card_lastname,
				cards.surname card_surname,
				cards.phone card_phone,
				schools.id school_id,
				schools.name school_name,
				books.id book_id,
				books.name book_name,
				books.author book_author,
				books.book_id book_ownID
			FROM
				cards
			LEFT JOIN schools
				ON schools.id = cards.school_id
			LEFT JOIN cards_books
				ON cards_books.card_id = cards.id
			LEFT JOIN books
				ON books.id = cards_books.book_id`,
}
