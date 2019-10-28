# Auto Fan
Auto-Fan let you the fan working with a relay and a temperature sensor together.
The temperature sensor will trigger the relay to control the fan running or stopping.

## Connect
temperature sensor:
- vcc: phys.1/3.3v
- dat: phys.7/BCM.4
- gnd: phys.9/GND

realy:
- vcc: phys.2/5v
- in:  phys.26/BCM.7
- gnd: phys.34/GND
- on:  fan(+)
- com: bettery(+)

```go

           temperature
           sensor
          o---------o
          |         |
          |         |
          |         |                                           o---------------o
          |         |               +-----------+               |fan            |
          |         |     +---------+ * 1   2 * +---------+     |    \_   _/    |
          o-+--+--+-o     |         | o       o |         |     |      \ /      |
            |  |  |       |  +------+ * 4     o |         |     |   ----*----   |    +----------+    +----------+
          gnd dat vcc     |  |  +---+ * 9     o |         |     |     _/ \_     |    |          |    |          |
            |  |  +-------+  |  |   | o       o |         |     |    /     \    |    |          |    |          |
            |  |             |  |   | o       o |         |     |               |    |       o-----------o      |
            |  +-------------+  |   | o       o |         |     o-----+---+-----+    |       |  -    +   |      |
            |                   |   | o       o |         |           |   |          |       |           |      |
            +-------------------+   | o       o |         |           |   |          |       |   power   |      |
                                    | o       o |         |           |   +----------+       o-----------o      |
                                    | o       o |         |           |                                         |
                                    | o       o |         |           |    +------------------------------------+
                                    | o    26 * +------+  |           NO   COM
                                    | o       o |      |  |           |    |
                                    | o       o |      |  |     o-----+----+---o
                                    | o       o |      |  +-vcc-+        relay |
                                    | o    34 * +--+   |        |    /         |
                                    | o       o |  +--------gnd-+   /          |
                                    | o       o |      |        |  o  ------o  |
                                    | o 39 40 o |      +-----in-+              |
                                    +-----------+               o--------------o

```

