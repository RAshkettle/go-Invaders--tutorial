package main

const (
	NUMBER_OF_ALIENS_IN_ROW = 12
	ALIEN_SIZE              = 16
	PADDING                 = 64
)

func SpawnAlienWave() []*Alien {
	aliens := make([]*Alien, 0)

	for i := range NUMBER_OF_ALIENS_IN_ROW {
		//Make the top row dude
		alien := NewAlien(SquidAlien)
		alien.X = (i * ALIEN_SIZE) + PADDING
		alien.Y = ALIEN_SIZE
		aliens = append(aliens, alien)

		//Make the Middle Row Dudes
		alien = NewAlien(ArmAlien)
		alien.X = (i * ALIEN_SIZE) + PADDING
		alien.Y = ALIEN_SIZE * 2
		aliens = append(aliens, alien)

		alien = NewAlien(ArmAlien)
		alien.X = (i * ALIEN_SIZE) + PADDING
		alien.Y = ALIEN_SIZE * 3
		aliens = append(aliens, alien)

		alien = NewAlien(FootAlien)
		alien.X = (i * ALIEN_SIZE) + PADDING
		alien.Y = ALIEN_SIZE * 4
		aliens = append(aliens, alien)

		alien = NewAlien(FootAlien)
		alien.X = (i * ALIEN_SIZE) + PADDING
		alien.Y = ALIEN_SIZE * 5
		aliens = append(aliens, alien)
	}
	return aliens
}
