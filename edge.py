import sqlite3


def main():
    # Open database connection
    with sqlite3.connect('edge.db') as conn:

        # Get database cursor
        cursor = conn.cursor()

        # Create table(s)
        cursor.execute('''
            CREATE TABLE IF NOT EXISTS history (
                speed INTEGER,
                heading INTEGER,
                temperature INTEGER,
                longitude REAL,
                latitude REAL
            );
        ''')

        # Insert a row of data
        cursor.execute("INSERT INTO history (speed,heading,temperature,longitude,latitude) VALUES (15,345,99,-71.0000,25.0000)")

        # Save (commit) the changes
        conn.commit()


if __name__ == '__main__':
    main()
