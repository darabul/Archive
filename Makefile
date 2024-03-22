DEFAULT_DIR_TO_ARCHIVE = ../test
DEFAULT_DIR_FOR_ARCHIVE = ../archive

DEFAULT_DIR_TO_UNARCHIVE = ../archive
DEFAULT_DIR_FOR_UNARCHIVE = ../unarchive

arch:
	go run cmd/main.go -m arch -s ${DEFAULT_DIR_TO_ARCHIVE} -r ${DEFAULT_DIR_FOR_ARCHIVE}

unarch:
	go run cmd/main.go -m unarch -s ${DEFAULT_DIR_TO_UNARCHIVE} -r ${DEFAULT_DIR_FOR_UNARCHIVE}