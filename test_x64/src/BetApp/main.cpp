#include "BetDialog.h"
#include "BetApplication.h"


int main(int argc, char *argv[])
{
	BetApplication a(argc, argv);

	BetDialog w;
    w.show();

    return a.exec();
}
