<img src="../../img/auto-fan.gif" width=40% height=40% />

# Auto Fan
Auto-Fan let you the fan working with a relay and a temperature sensor together.
The temperature sensor will trigger the relay to control the fan running or stopping.

## Connect
temperature sensor:
- vcc: 3.3v
- dat: GPIO-4
- gnd: GND

realy:
- vcc: 5v
- in:  GPIO-7
- gnd: GND
- on:  fan(+)
- com: bettery(+)

```go

 o---------o
 |         |
 |temperature
 |sensor   |                                           o---------------o
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
