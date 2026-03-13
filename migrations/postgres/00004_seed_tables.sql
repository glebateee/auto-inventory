-- +goose Up
INSERT INTO categories (name) VALUES
('Brake System'),
('Engine Parts'),
('Suspension'),
('Electrical'),
('Exhaust System'),
('Transmission');

-- Insert manufacturers (4 records)
INSERT INTO manufacturers (name, country) VALUES
('Bosch', 'Germany'),
('Denso', 'Japan'),
('Valeo', 'France'),
('Delphi', 'USA');

-- Insert products (20 records)
INSERT INTO products (
    sku, name, description, category_id, manufacturer_id, weight, unit, price, baseprice, issueyear
) VALUES
-- Brake System (category_id = 1)
('BP-1234', 'Brake Pad Set', 'Ceramic brake pads for sedans', 1, 1, 2, 'pc.', 4500, 3200, 2023),
('BP-5678', 'Brake Disc', 'Vented front brake disc', 1, 2, 8, 'pc.', 7800, 5900, 2022),
('BP-9012', 'Brake Caliper', 'Front left brake caliper', 1, 3, 3, 'pc.', 12500, 9800, 2021),
('BP-3456', 'Brake Hose', 'Stainless steel braided brake hose', 1, 4, 1, 'pc.', 2100, 1500, 2023),

-- Engine Parts (category_id = 2)
('EP-7890', 'Oil Filter', 'High-performance oil filter', 2, 1, 1, 'pc.', 850, 600, 2024),
('EP-2345', 'Air Filter', 'Panel air filter', 2, 2, 1, 'pc.', 1200, 900, 2023),
('EP-6789', 'Spark Plug', 'Iridium spark plug (set of 4)', 2, 3, 1, 'set', 5500, 4200, 2022),
('EP-0123', 'Water Pump', 'Engine water pump with gasket', 2, 4, 2, 'pc.', 8900, 6700, 2021),

-- Suspension (category_id = 3)
('SP-4567', 'Shock Absorber', 'Front shock absorber', 3, 1, 4, 'pc.', 10200, 8100, 2023),
('SP-8901', 'Control Arm', 'Lower control arm with ball joint', 3, 2, 5, 'pc.', 13500, 10700, 2022),
('SP-2345', 'Ball Joint', 'Suspension ball joint', 3, 3, 2, 'pc.', 3200, 2400, 2024),
('SP-6789', 'Stabilizer Link', 'Front stabilizer link', 3, 4, 1, 'pc.', 1800, 1300, 2023),

-- Electrical (category_id = 4)
('EL-0123', 'Alternator', '12V 120A alternator', 4, 1, 6, 'pc.', 24500, 19200, 2022),
('EL-4567', 'Starter Motor', '1.4kW starter', 4, 2, 4, 'pc.', 18700, 14900, 2021),
('EL-8901', 'Ignition Coil', 'Ignition coil pack', 4, 3, 1, 'pc.', 4300, 3300, 2024),

-- Exhaust System (category_id = 5)
('EX-2345', 'Catalytic Converter', 'Universal catalytic converter', 5, 4, 5, 'pc.', 32800, 26500, 2020),
('EX-6789', 'Oxygen Sensor', 'Lambda sensor (wideband)', 5, 1, 1, 'pc.', 6700, 5200, 2023),

-- Transmission (category_id = 6)
('TR-1234', 'Clutch Kit', 'Clutch disc, pressure plate, release bearing', 6, 2, 9, 'set', 35600, 28900, 2022),
('TR-5678', 'Gearbox Mount', 'Manual transmission mount', 6, 3, 2, 'pc.', 4100, 3200, 2023),
('TR-9012', 'CV Joint', 'Outer CV joint', 6, 4, 2, 'pc.', 5900, 4500, 2024)
ON CONFLICT DO NOTHING ;

-- +goose Down
-- Delete products (by SKU)
DELETE FROM products WHERE sku IN (
    'BP-1234', 'BP-5678', 'BP-9012', 'BP-3456',
    'EP-7890', 'EP-2345', 'EP-6789', 'EP-0123',
    'SP-4567', 'SP-8901', 'SP-2345', 'SP-6789',
    'EL-0123', 'EL-4567', 'EL-8901',
    'EX-2345', 'EX-6789',
    'TR-1234', 'TR-5678', 'TR-9012'
);

-- Delete manufacturers (by ID)
DELETE FROM manufacturers WHERE id IN (1, 2, 3, 4);

-- Delete categories (by ID)
DELETE FROM categories WHERE id IN (1, 2, 3, 4, 5, 6);
